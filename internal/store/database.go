package store

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func OpenPool() (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return dbPool, nil
}

func MigrateFS(dbPool *pgxpool.Pool, migrationFs fs.FS, dir string) error {
	goose.SetBaseFS(migrationFs)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(dbPool, dir)
}

func Migrate(dbPool *pgxpool.Pool, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migration dialect setup failed: %w", err)
	}

	db := stdlib.OpenDBFromPool(dbPool)
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
