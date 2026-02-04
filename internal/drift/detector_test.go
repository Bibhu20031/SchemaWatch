package drift

import (
	"testing"

	"github.com/Bibhu20031/SchemaWatch/internal/snapshot"
)

func TestDetectSchemaDrift(t *testing.T) {
	prev := []snapshot.Column{
		{
			Name:     "id",
			DataType: "integer",
			Nullable: false,
		},
		{
			Name:     "email",
			DataType: "text",
			Nullable: false,
		},
	}

	defaultNow := "now()"

	curr := []snapshot.Column{
		{
			Name:     "id",
			DataType: "integer",
			Nullable: false,
		},
		{
			Name:     "email",
			DataType: "text",
			Nullable: true, // changed
		},
		{
			Name:       "created_at",
			DataType:   "timestamp",
			Nullable:   false,
			DefaultVal: &defaultNow, // added
		},
	}

	changes := Detect(prev, curr)

	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(changes))
	}

	foundNullable := false
	foundAdded := false

	for _, c := range changes {
		if c.Type == NullabilityChanged && c.ColumnName == "email" {
			foundNullable = true
		}
		if c.Type == ColumnAdded && c.ColumnName == "created_at" {
			foundAdded = true
		}
	}

	if !foundNullable || !foundAdded {
		t.Fatalf("unexpected drift changes: %+v", changes)
	}
}
