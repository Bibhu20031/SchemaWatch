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
