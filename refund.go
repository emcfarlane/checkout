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

func (s *Service) Refund(ctx context.Context, req *pb.RefundRequest) (*pb.Authorization, error) {
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

	// Can only refund if capturing transaction or already refunding.
	if state != pb.Authorization_CAPTURE && state != pb.Authorization_REFUND {
		tx.Rollback()
		return nil, status.Errorf(codes.FailedPrecondition, "authorization %s can't be refunded", state)
	}

	if req.Amount > amountCaptured {
		tx.Rollback()
		return nil, status.Errorf(codes.InvalidArgument, "refund amount %d greater than captured %d", req.Amount, amountCaptured)
	}
	amountCaptured -= req.Amount

	state = pb.Authorization_REFUND
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
	s.log.Info("refunded authorization", zap.String("id", id), zap.Uint64("amount", req.Amount))

	return &pb.Authorization{
		Id:             id,
		State:          state,
		Amount:         amount,
		AmountCaptured: amountCaptured,
		CreateTime:     timestamppb.New(createTime),
		UpdateTime:     timestamppb.New(updateTime),
	}, nil
}
