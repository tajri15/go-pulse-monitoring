package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	conn *pgxpool.Pool
}

func NewStore(conn *pgxpool.Pool) *Store {
	return &Store{conn: conn}
}

// Definisikan struct User agar sesuai dengan tabel di database
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Jangan kirim hash password ke client
	CreatedAt    time.Time `json:"created_at"`
}

type CreateUserParams struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

// Method untuk membuat user baru
func (s *Store) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, username, email, password_hash, created_at`
	
	row := s.conn.QueryRow(ctx, query, arg.Username, arg.Email, arg.PasswordHash)
	
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)
	return u, err
}

// Method untuk mendapatkan user berdasarkan email
func (s *Store) GetUserByEmail(ctx context.Context, email string) (User, error) {
	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE email = $1 LIMIT 1`
	
	row := s.conn.QueryRow(ctx, query, email)

	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.CreatedAt)
	return u, err
}