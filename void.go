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

func (s *Service) Void(ctx context.Context, req *pb.VoidRequest) (*pb.Authorization, error) {
	id := req.Id
	updateTime := time.Now()

	stmt, args, err := s.psql.Update("authorizations").
		Set("state", pb.Authorization_VOID).
		Set("update_time", updateTime).
		Where(sq.Eq{"id": id}).
		Where(sq.Eq{"state": pb.Authorization_OPEN}).
		Suffix("RETURNING amount, amount_captured, create_time").
		ToSql()
	if err != nil {
		return nil, err
	}
	s.log.Debug("void authorization", zap.String("sql", stmt))

	// Void the authorization.
	var amount, amountCaptured uint64
	var createTime time.Time
	if err := s.db.QueryRowContext(ctx, stmt, args...).Scan(&amount, &amountCaptured, &createTime); err != nil {
		if err != sql.ErrNoRows {
			return nil, err // unknown error
		}

		// Query to see why update failed.
		var state pb.Authorization_State
		if err := s.psql.Select("state").
			From("authorizations").
			Where(sq.Eq{"id": id}).
			RunWith(s.db).
			QueryRowContext(ctx).
			Scan(&state); err != nil {
			if err != sql.ErrNoRows {
				return nil, err
			}
			return nil, status.Errorf(codes.NotFound, "authorization doesn't exist")
		}
		// Update failed due to invalid state.
		if state == pb.Authorization_VOID {
			return nil, status.Errorf(codes.FailedPrecondition, "authorization already voided")
		}
		return nil, status.Errorf(codes.FailedPrecondition, "authorization %s cannot be voided", state)
	}
	s.log.Info("voided authorization", zap.String("id", id))

	return &pb.Authorization{
		Id:             id,
		State:          pb.Authorization_VOID,
		Amount:         amount,
		AmountCaptured: amountCaptured,
		CreateTime:     timestamppb.New(createTime),
		UpdateTime:     timestamppb.New(updateTime),
	}, nil
}
