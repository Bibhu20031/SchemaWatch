package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) IsValidKey(ctx context.Context, keyHash string) (bool, error) {
	var exists bool

	err := r.db.QueryRow(
		ctx,
		`SELECT EXISTS (
			SELECT 1 FROM api_keys
			WHERE key_hash = $1 AND is_active = true
		)`,
		keyHash,
	).Scan(&exists)

	return exists, err
}
