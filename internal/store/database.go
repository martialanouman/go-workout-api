package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=workout_user password=workout_password dbname=workout_db sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("db: open %w", err)
	}

	fmt.Println("Connected to database successfully!")
	return db, nil
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}

func MigrateFS(db *sql.DB, migrationFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationFS)

	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}
