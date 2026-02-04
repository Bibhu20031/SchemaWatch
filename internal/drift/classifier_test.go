package drift

import "testing"

func TestClassifyBreakingChange(t *testing.T) {
	changes := []DriftChange{
		{
			Type:        ColumnRemoved,
			ColumnName:  "email",
			BeforeValue: "text",
		},
	}

	classified := Classify(changes)

	if len(classified) != 1 {
		t.Fatalf("expected 1 change, got %d", len(classified))
	}

	if classified[0].Impact != ImpactBreaking {
		t.Fatalf("expected BREAKING, got %s", classified[0].Impact)
	}
}
