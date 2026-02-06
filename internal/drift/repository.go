package drift

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

func (r *Repository) StoreEvents(
	ctx context.Context,
	schemaID int64,
	fromVersion int,
	toVersion int,
	changes []ClassifiedChange,
) error {

	for _, c := range changes {
		var beforeJSON, afterJSON []byte

		if c.BeforeValue != nil {
			beforeJSON, _ = json.Marshal(c.BeforeValue)
		}
		if c.AfterValue != nil {
			afterJSON, _ = json.Marshal(c.AfterValue)
		}

		_, err := r.db.Exec(ctx, `
			INSERT INTO drift_events (
				schema_id,
				version_from,
				version_to,
				change_type,
				column_name,
				impact,
				before_value,
				after_value
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		`,
			schemaID,
			fromVersion,
			toVersion,
			c.Type,
			c.ColumnName,
			c.Impact,
			beforeJSON,
			afterJSON,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) ListBySchema(
	ctx context.Context,
	schemaID int64,
) ([]map[string]any, error) {

	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			version_from,
			version_to,
			change_type,
			column_name,
			impact,
			before_value,
			after_value,
			created_at
		FROM drift_events
		WHERE schema_id = $1
		ORDER BY created_at DESC
	`, schemaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]any

	for rows.Next() {
		var (
			id          int64
			fromV       int
			toV         int
			changeType  string
			columnName  string
			impact      string
			beforeValue []byte
			afterValue  []byte
			createdAt   any
		)

		if err := rows.Scan(
			&id,
			&fromV,
			&toV,
			&changeType,
			&columnName,
			&impact,
			&beforeValue,
			&afterValue,
			&createdAt,
		); err != nil {
			return nil, err
		}

		result = append(result, map[string]any{
			"id":           id,
			"version_from": fromV,
			"version_to":   toV,
			"change_type":  changeType,
			"column_name":  columnName,
			"impact":       impact,
			"before_value": beforeValue,
			"after_value":  afterValue,
			"created_at":   createdAt,
		})
	}

	return result, nil
}
