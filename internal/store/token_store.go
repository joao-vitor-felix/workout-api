package store

import (
	"database/sql"
	"time"

	"github.com/joao-vitor-felix/workout-api/internal/tokens"
)

type TokenStore interface {
	Insert(token *tokens.Token) error
	Create(userId int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteForUser(userId int, scope string) error
}

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db}
}

func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, scope, expires_at)
		VALUES ($1, $2, $3, $4)`
	_, err := t.db.Exec(query, token.Hash, token.UserID, token.Scope, token.ExpiresAt)
	return err
}

func (t *PostgresTokenStore) Create(userId int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = t.Insert(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (t *PostgresTokenStore) DeleteForUser(userId int, scope string) error {
	query := `
		DELETE FROM tokens
		WHERE user_id = $1 AND scope = $2`
	_, err := t.db.Exec(query, userId, scope)
	return err
}
