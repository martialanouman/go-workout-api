package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainText string) error {
	const cost = 12
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), cost)
	if err != nil {
		return err
	}

	p.plainText = &plainText
	p.hash = hash

	return nil
}

func (p *password) Matches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type User struct {
	Id           int64     `json:"id"`
	Username     string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db}
}

type UserStore interface {
	CreateUser(*User) error
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
	GetUserByToken(scope, tokenPlaintext string) (*User, error)
}

func (p *PostgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users (username, email, password_hash, bio)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
	`

	err := p.db.QueryRow(
		query, user.Username, user.Email, user.PasswordHash.hash, user.Bio,
	).Scan(
		&user.Id, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}

	query := `
	SELECT id, username, email, password_hash, bio, created_at, updated_at
	FROM users
	WHERE username = $1
	`

	err := p.db.QueryRow(query, username).Scan(
		&user.Id, &user.Username, &user.Email, &user.PasswordHash.hash,
		&user.Bio, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *PostgresUserStore) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET username=$1, email=$2, bio=$3, updated_at=CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at
	`

	result, err := p.db.Exec(query, user.Username, user.Email, user.Bio, user.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (p *PostgresUserStore) GetUserByToken(scope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	user := &User{
		PasswordHash: password{},
	}

	query := `
	SELECT u.id, u.username, u.email, u.password_hash, u.bio, u.created_at, u.updated_at
	FROM users u
	INNER JOIN tokens t ON t.user_id = u.id
	WHERE t.hash = $1 AND scope = $2 AND t.expiry > $3
	`

	err := p.db.QueryRow(query, tokenHash[:], scope, time.Now()).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
