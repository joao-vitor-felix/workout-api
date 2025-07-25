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
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.plainText = &plainText
	p.hash = hash
	return nil
}

func (p *password) Check(plainText string) (bool, error) {
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
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash password  `json:"-"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type UserStore interface {
	Create(user *User) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(*User) (*User, error)
	GetUserToken(scope, tokenPlainText string) (*User, error)
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db}
}

func (s *PostgresUserStore) Create(user *User) (*User, error) {
	query := `
    INSERT INTO users (email, username, password_hash, bio)
    VALUES ($1, $2, $3, $4)
    RETURNING id, created_at, updated_at
  `

	err := s.db.QueryRow(query, user.Email, user.Username, user.PasswordHash.hash, user.Bio).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) GetByUsername(username string) (*User, error) {
	user := &User{
		Username:     username,
		PasswordHash: password{},
	}

	query := `
  SELECT id, email, password_hash, bio, created_at, updated_at
  FROM users
  WHERE username = $1
  `

	err := s.db.QueryRow(query, username).Scan(&user.ID, &user.Email, &user.PasswordHash.hash, &user.Bio, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) Update(user *User) (*User, error) {
	query := `
    UPDATE users
    SET email = $1, username = $2, password_hash = $3, bio = $4, updated_at = NOW()
    WHERE id = $5
    RETURNING updated_at
  `

	err := s.db.QueryRow(query, user.Email, user.Username, user.PasswordHash.hash, user.Bio, user.ID).Scan(&user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) GetUserToken(scope, plaintextPassword string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(plaintextPassword))

	query := `
  SELECT
    u.id,
    u.username,
    u.email,
    u.password_hash,
    u.bio,
    u.created_at,
    u.updated_at
  FROM
    users u
  INNER JOIN
    tokens t ON t.user_id = u.id
  WHERE
    t.hash = $1
  AND
    t.scope = $2 and t.expires_at > $3
  `

	user := &User{
		PasswordHash: password{},
	}

	err := s.db.QueryRow(query, tokenHash[:], scope, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
