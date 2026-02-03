package schema

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateSchema(
	ctx context.Context,
	req RegisterRequest,
	snapshot any,
) (int64, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var schemaID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO schemas (db_host, db_port, db_name, schema_name, table_name)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id
	`,
		req.DBHost, req.DBPort, req.DBName, req.SchemaName, req.TableName,
	).Scan(&schemaID)
	if err != nil {
		return 0, err
	}

	data, _ := json.Marshal(snapshot)

	_, err = tx.Exec(ctx, `
		INSERT INTO schema_versions (schema_id, version, snapshot)
		VALUES ($1, 1, $2)
	`, schemaID, data)
	if err != nil {
		return 0, err
	}

	return schemaID, tx.Commit(ctx)
}
