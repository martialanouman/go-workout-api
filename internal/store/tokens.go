package store

import (
	"database/sql"
	"time"

	"github.com/martialanouman/femProject/internal/tokens"
)

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateToken(userId int64, ttl time.Duration, scope string) (*tokens.Token, error)
	RevokeAllTokenForUser(userId int, scope string) error
}

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db}
}

func (s *PostgresTokenStore) CreateToken(userId int64, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = s.Insert(token)
	return token, err
}

func (s *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)
	`

	_, err := s.db.Exec(query, token.Hash, token.UserId, token.Expiry, token.Scope)

	return err

}

func (s *PostgresTokenStore) RevokeAllTokenForUser(userId int, scope string) error {
	query := `
		DELETE FROM tokens
		WHERE user_id = $1 AND scope = $2
	`

	_, err := s.db.Exec(query, userId, scope)

	return err
}
