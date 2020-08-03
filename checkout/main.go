package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"

	"github.com/emcfarlane/checkout"
	"go.uber.org/zap"
	"gocloud.dev/postgres"
)

func env(key, def string) string {
	if e := os.Getenv(key); e != "" {
		return e
	}
	return def
}

var (
	flagAddress  = flag.String("addr", env("ADDRESS", ":8080"), "Address to listen on")
	flagDatabase = flag.String("pg", env("POSTGRES", "postgres://edward:password@localhost/checkout"), "Postgres database URL")
)

func run() error {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zap.NewExample()

	db, err := postgres.Open(ctx, *flagDatabase)
	if err != nil {
		return err
	}
	defer db.Close()

	// Mock bank connection for routing payments.
	bank := checkout.NewMockBank()

	s, err := checkout.NewService(logger, db, bank)
	if err != nil {
		return err
	}
	if err := s.CreateTables(ctx); err != nil {
		return err
	}

	l, err := net.Listen("tcp", *flagAddress)
	if err != nil {
		return err
	}

	return s.Serve(l)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
