// Package dbtest gives integration tests an isolated PostgreSQL schema with every migration applied.
package dbtest

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// Open returns a database whose search_path is a fresh schema. Tests are skipped unless TRANSAT_TEST_DSN is set,
// e.g. "host=/tmp port=54329 user=postgres dbname=transat_test sslmode=disable". The database name must contain
// "test" so a real database is never touched.
func Open(t testing.TB) *sql.DB {
	t.Helper()
	dsn := os.Getenv("TRANSAT_TEST_DSN")
	if dsn == "" {
		t.Skip("TRANSAT_TEST_DSN not set")
	}
	if !strings.Contains(dsn, "test") {
		t.Fatal("TRANSAT_TEST_DSN must target a database whose name contains \"test\"")
	}

	admin, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	schema := fmt.Sprintf("t_%d_%d", os.Getpid(), time.Now().UnixNano())
	if _, err := admin.Exec("CREATE SCHEMA " + schema); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = admin.Exec("DROP SCHEMA " + schema + " CASCADE")
		_ = admin.Close()
	})

	connector, err := pq.NewConnector(dsn + " search_path=" + schema)
	if err != nil {
		t.Fatal(err)
	}
	db := sql.OpenDB(connector)
	db.SetMaxOpenConns(16)
	t.Cleanup(func() { _ = db.Close() })

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(file), "..", "..", "..", "db", "migrations")
	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatal(err)
	}
	goose.SetLogger(goose.NopLogger())
	if err := goose.Up(db, dir); err != nil {
		t.Fatalf("migrations: %v", err)
	}
	return db
}
