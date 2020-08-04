package checkout

import (
	"context"
	"database/sql"
	"net"
	"net/http"

	sq "github.com/Masterminds/squirrel"
	"github.com/emcfarlane/checkout/checkoutpb"
	"github.com/emcfarlane/graphpb"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service implements a CheckoutServer.
type Service struct {
	checkoutpb.UnimplementedCheckoutServer

	log  *zap.Logger
	db   *sql.DB
	bank Bank
	psql sq.StatementBuilderType // Postgres builder
}

func NewService(logger *zap.Logger, db *sql.DB, bank Bank) (*Service, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &Service{
		log:  logger,
		db:   db,
		bank: bank,
		psql: psql,
	}, nil
}

func createTables(db *sql.DB, ctx context.Context) error {
	_, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS authorizations (
	id UUID PRIMARY KEY,
	state INTEGER NOT NULL,
	amount INTEGER NOT NULL,
	amount_captured INTEGER NOT NULL,
	create_time TIMESTAMP NOT NULL,
	update_time TIMESTAMP NOT NULL
);`)
	return err
}

// CreateTables initialises the postgres database.
func (s *Service) CreateTables(ctx context.Context) error {
	if err := createTables(s.db, ctx); err != nil {
		s.log.Error("database init error", zap.Error(err))
		return err
	}
	return nil
}

// unaryInterceptor sits inbetween every gRPC call.
func (s *Service) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// TODO: authentication could be implemented as part of this interceptor.

	resp, err = handler(ctx, req)
	if err != nil {
		if ss, ok := status.FromError(err); ok {
			s.log.Error(ss.Message(), zap.String("code", ss.Code().String()))
		} else {
			s.log.Error(err.Error(), zap.String("code", codes.Internal.String()))
		}
	}
	return
}

// Serve blocks serving on the listener a gRPC and REST based API.
func (s *Service) Serve(l net.Listener) error {
	m := cmux.New(l)

	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.Any())

	zapgrpc.NewLogger(s.log)

	gs := grpc.NewServer(
		grpc.UnaryInterceptor(s.unaryInterceptor),
	)
	checkoutpb.RegisterCheckoutServer(gs, s)

	// graphpb implements the REST based HTTP bindings documented in the
	// protobuf file checkoutpb/checkout.proto.
	hd := &graphpb.Handler{
		UnaryInterceptor: s.unaryInterceptor,
	}
	hd.RegisterServiceByName("checkout.api.Checkout", s)

	hs := &http.Server{
		Handler: hd,
	}

	errs := make(chan error)

	go func() { errs <- gs.Serve(grpcL) }()
	defer gs.Stop()

	go func() { errs <- hs.Serve(httpL) }()
	defer hs.Close()

	go func() { errs <- m.Serve() }()

	s.log.Info("listening", zap.String("address", l.Addr().String()))
	select {
	case err := <-errs:
		return err
	}
}
