// +build integration

package checkout

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"gocloud.dev/postgres"
)

var db *sql.DB // Test database connection

func env(key, def string) string {
	if e := os.Getenv(key); e != "" {
		return e
	}
	return def
}

func TestMain(m *testing.M) {
	dbURL := env("POSTGRES", "postgres://edward:password@localhost/checkout")
	ctx := context.Background()

	var err error
	db, err = postgres.Open(ctx, dbURL)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialise tables.
	if err := createTables(db, ctx); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
