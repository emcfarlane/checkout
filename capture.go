package checkout

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	pb "github.com/emcfarlane/checkout/checkoutpb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) Capture(ctx context.Context, req *pb.CaptureRequest) (*pb.Authorization, error) {
	id := req.Id
	if req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be greater than zero")
	}
	updateTime := time.Now()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var state pb.Authorization_State
	var amount, amountCaptured uint64
	var createTime time.Time

	if err := s.psql.Select("state", "amount", "amount_captured", "create_time").
		From("authorizations").
		Where(sq.Eq{"id": id}).
		RunWith(tx).
		QueryRowContext(ctx).
		Scan(&state, &amount, &amountCaptured, &createTime); err != nil {

		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "authorization doesn't exist")
		}
		return nil, err
	}

	// Can only capture if newly open transaction or already capturing.
	if state != pb.Authorization_OPEN && state != pb.Authorization_CAPTURE {
		tx.Rollback()
		return nil, status.Errorf(codes.FailedPrecondition, "authorization %s can't be captured", state)
	}

	if req.Amount > (amount - amountCaptured) {
		tx.Rollback()
		return nil, status.Errorf(codes.InvalidArgument, "capture amount %d greater than %d authorized", req.Amount, amount)
	}

	// Capture funds with the Bank.
	if err := s.bank.Capture(id, req.Amount); err != nil {
		tx.Rollback()
		return nil, err
	}
	amountCaptured += req.Amount

	state = pb.Authorization_CAPTURE
	if _, err := s.psql.Update("authorizations").
		Set("state", state).
		Set("amount_captured", amountCaptured).
		Set("update_time", updateTime).
		Where(sq.Eq{"id": id}).
		RunWith(tx).
		ExecContext(ctx); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	s.log.Info("captured authorization", zap.String("id", id), zap.Uint64("amount", req.Amount))

	return &pb.Authorization{
		Id:             id,
		State:          state,
		Amount:         amount,
		AmountCaptured: amountCaptured,
		CreateTime:     timestamppb.New(createTime),
		UpdateTime:     timestamppb.New(updateTime),
	}, nil

}
