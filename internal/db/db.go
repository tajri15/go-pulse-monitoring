package db

import (
	"context"
	"errors"
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

// Definisikan struct Site
type Site struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateSiteParams struct {
	UserID int64  `json:"user_id"`
	URL    string `json:"url"`
}

// Method untuk membuat site baru
func (s *Store) CreateSite(ctx context.Context, arg CreateSiteParams) (Site, error) {
	query := `INSERT INTO sites (user_id, url) VALUES ($1, $2) RETURNING id, user_id, url, created_at`
	
	row := s.conn.QueryRow(ctx, query, arg.UserID, arg.URL)

	var site Site
	err := row.Scan(&site.ID, &site.UserID, &site.URL, &site.CreatedAt)
	return site, err
}

// Method untuk mendapatkan semua site milik seorang user
func (s *Store) GetSitesByUserID(ctx context.Context, userID int64) ([]Site, error) {
	query := `SELECT id, user_id, url, created_at FROM sites WHERE user_id = $1 ORDER BY created_at DESC`
	
	rows, err := s.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sites := []Site{}
	for rows.Next() {
		var site Site
		if err := rows.Scan(&site.ID, &site.UserID, &site.URL, &site.CreatedAt); err != nil {
			return nil, err
		}
		sites = append(sites, site)
	}
	return sites, nil
}

// Method untuk menghapus site, dengan verifikasi kepemilikan
func (s *Store) DeleteSite(ctx context.Context, siteID int64, userID int64) error {
	query := `DELETE FROM sites WHERE id = $1 AND user_id = $2`
	
	cmdTag, err := s.conn.Exec(ctx, query, siteID, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("site not found or user not authorized to delete")
	}
	return nil
}
