package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		panic("error loading .env file")
	}

	var (
		HOST = os.Getenv("DB_HOST")
		PORT = os.Getenv("DB_PORT")
		USER = os.Getenv("DB_USER")
		PASS = os.Getenv("DB_PASSWORD")
		DB   = os.Getenv("DB_NAME")
	)

	db, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, USER, PASS, DB))
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
