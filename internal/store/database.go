package store

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

func OpenPool() (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return dbPool, nil
}

func MigrateFS(db *sql.DB, migrationFs fs.FS, dir string) error {
	goose.SetBaseFS(migrationFs)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migration dialect setup failed: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
