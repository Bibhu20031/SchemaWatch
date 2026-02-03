package snapshot

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Column struct {
	Name       string  `json:"name"`
	DataType   string  `json:"data_type"`
	Nullable   bool    `json:"nullable"`
	DefaultVal *string `json:"default,omitempty"`
}

func FetchTableSchema(
	ctx context.Context,
	pool *pgxpool.Pool,
	schema string,
	table string,
) ([]Column, error) {

	rows, err := pool.Query(ctx, `
		SELECT column_name, data_type, is_nullable, column_default
		FROM information_schema.columns
		WHERE table_schema = $1 AND table_name = $2
		ORDER BY ordinal_position
	`, schema, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []Column
	for rows.Next() {
		var c Column
		var nullable string
		err := rows.Scan(&c.Name, &c.DataType, &nullable, &c.DefaultVal)
		if err != nil {
			return nil, err
		}
		c.Nullable = nullable == "YES"
		cols = append(cols, c)
	}

	return cols, nil
}
