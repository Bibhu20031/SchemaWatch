package drift

type ChangeType string

const (
	ColumnAdded        ChangeType = "COLUMN_ADDED"
	ColumnRemoved      ChangeType = "COLUMN_REMOVED"
	DataTypeChanged    ChangeType = "DATA_TYPE_CHANGED"
	NullabilityChanged ChangeType = "NULLABILITY_CHANGED"
	DefaultChanged     ChangeType = "DEFAULT_CHANGED"
)

type DriftChange struct {
	Type        ChangeType `json:"type"`
	ColumnName  string     `json:"column_name"`
	BeforeValue any        `json:"before,omitempty"`
	AfterValue  any        `json:"after,omitempty"`
}

type ImpactLevel string

const (
	ImpactSafe     ImpactLevel = "SAFE"
	ImpactRisky    ImpactLevel = "RISKY"
	ImpactBreaking ImpactLevel = "BREAKING"
)

type ClassifiedChange struct {
	DriftChange
	Impact ImpactLevel `json:"impact"`
}
