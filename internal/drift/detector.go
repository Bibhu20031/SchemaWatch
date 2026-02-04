package drift

import (
	"github.com/Bibhu20031/SchemaWatch/internal/snapshot"
)

func Detect(
	prev []snapshot.Column,
	curr []snapshot.Column,
) []DriftChange {

	changes := []DriftChange{}

	prevMap := make(map[string]snapshot.Column)
	currMap := make(map[string]snapshot.Column)

	for _, c := range prev {
		prevMap[c.Name] = c
	}
	for _, c := range curr {
		currMap[c.Name] = c
	}

	// Added columns
	for name, col := range currMap {
		if _, ok := prevMap[name]; !ok {
			changes = append(changes, DriftChange{
				Type:       ColumnAdded,
				ColumnName: name,
				AfterValue: col,
			})
		}
	}

	// Removed columns
	for name, col := range prevMap {
		if _, ok := currMap[name]; !ok {
			changes = append(changes, DriftChange{
				Type:        ColumnRemoved,
				ColumnName:  name,
				BeforeValue: col,
			})
		}
	}

	// Modified columns
	for name, prevCol := range prevMap {
		currCol, ok := currMap[name]
		if !ok {
			continue
		}

		if prevCol.DataType != currCol.DataType {
			changes = append(changes, DriftChange{
				Type:        DataTypeChanged,
				ColumnName:  name,
				BeforeValue: prevCol.DataType,
				AfterValue:  currCol.DataType,
			})
		}

		if prevCol.Nullable != currCol.Nullable {
			changes = append(changes, DriftChange{
				Type:        NullabilityChanged,
				ColumnName:  name,
				BeforeValue: prevCol.Nullable,
				AfterValue:  currCol.Nullable,
			})
		}

		if !equalDefault(prevCol.DefaultVal, currCol.DefaultVal) {
			changes = append(changes, DriftChange{
				Type:        DefaultChanged,
				ColumnName:  name,
				BeforeValue: prevCol.DefaultVal,
				AfterValue:  currCol.DefaultVal,
			})
		}
	}

	return changes
}

func equalDefault(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
