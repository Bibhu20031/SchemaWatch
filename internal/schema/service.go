package schema

import (
	"context"

	"github.com/Bibhu20031/SchemaWatch/internal/snapshot"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repo *Repository
	db   *pgxpool.Pool
}

func NewService(repo *Repository, db *pgxpool.Pool) *Service {
	return &Service{repo: repo, db: db}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (int64, error) {
	snap, err := snapshot.FetchTableSchema(
		ctx,
		s.db,
		req.SchemaName,
		req.TableName,
	)
	if err != nil {
		return 0, err
	}

	return s.repo.CreateSchema(ctx, req, snap)
}

func (s *Service) List(ctx context.Context) ([]map[string]any, error) {
	return s.repo.ListSchemas(ctx)
}
