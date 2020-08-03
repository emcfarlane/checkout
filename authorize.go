package checkout

import (
	"context"
	"time"

	"github.com/emcfarlane/checkout/checkoutpb"
	pb "github.com/emcfarlane/checkout/checkoutpb"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*pb.Authorization, error) {
	pan, err := parsePAN(req.Pan)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid pan: %v", err.Error())
	}
	if ok := checkLuhn(pan); !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid pan: luhn check failure")
	}
	expiry, err := parseExpiry(req.ExpYear, req.ExpMonth)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid expiry: %v", err)
	}
	cvv := req.Cvv
	if ok := checkCVV(cvv); !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid cvv")
	}
	amount := req.Amount
	if amount <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid amount %d", req.Amount)
	}
	currency, err := parseCurrency(req.Currency)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id := uuid.New()
	t := time.Now()

	// Authorize the payment with the Bank.
	if err := s.bank.Authorize(id.String(), pan, cvv, expiry, amount, currency); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}

	// Create a new Authorization.
	sql, args, err := s.psql.Insert("authorizations").
		Columns("id", "state", "amount", "amount_captured", "create_time", "update_time").
		Values(id.String(), pb.Authorization_OPEN, amount, 0, t, t).
		ToSql()
	if err != nil {
		return nil, err
	}
	s.log.Debug("insert authorization", zap.String("sql", sql))

	if _, err := s.db.ExecContext(ctx, sql, args...); err != nil {
		return nil, err
	}
	s.log.Info("created authorization", zap.String("id", id.String()))

	ts := timestamppb.New(t)
	return &pb.Authorization{
		Id:             id.String(),
		State:          checkoutpb.Authorization_OPEN,
		Amount:         amount,
		AmountCaptured: 0,
		Currency:       currency,
		CreateTime:     ts,
		UpdateTime:     ts,
	}, nil
}
