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

func (r *Repository) ListSchemas(ctx context.Context) ([]map[string]any, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, db_host, db_port, db_name, schema_name, table_name, created_at
		FROM schemas
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]any

	for rows.Next() {
		var (
			id         int64
			dbHost     string
			dbPort     int
			dbName     string
			schemaName string
			tableName  string
			createdAt  any
		)

		if err := rows.Scan(
			&id, &dbHost, &dbPort, &dbName,
			&schemaName, &tableName, &createdAt,
		); err != nil {
			return nil, err
		}

		result = append(result, map[string]any{
			"id":          id,
			"db_host":     dbHost,
			"db_port":     dbPort,
			"db_name":     dbName,
			"schema_name": schemaName,
			"table_name":  tableName,
			"created_at":  createdAt,
		})
	}

	return result, nil
}

func (r *Repository) GetLatestVersion(
	ctx context.Context,
	schemaID int64,
) (int, []byte, error) {

	var version int
	var snapshot []byte

	err := r.db.QueryRow(ctx, `
		SELECT version, snapshot
		FROM schema_versions
		WHERE schema_id = $1
		ORDER BY version DESC
		LIMIT 1
	`, schemaID).Scan(&version, &snapshot)

	return version, snapshot, err
}
