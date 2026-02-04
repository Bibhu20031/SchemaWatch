package drift

func Classify(changes []DriftChange) []ClassifiedChange {
	result := make([]ClassifiedChange, 0, len(changes))

	for _, c := range changes {
		result = append(result, ClassifiedChange{
			DriftChange: c,
			Impact:      classifyOne(c),
		})
	}

	return result
}

func classifyOne(c DriftChange) ImpactLevel {
	switch c.Type {

	case ColumnAdded:
		return ImpactSafe

	case ColumnRemoved:
		return ImpactBreaking

	case DataTypeChanged:
		return ImpactBreaking

	case NullabilityChanged:

		if before, ok1 := c.BeforeValue.(bool); ok1 {
			if after, ok2 := c.AfterValue.(bool); ok2 {
				if before && !after {
					return ImpactBreaking
				}
			}
		}
		return ImpactRisky

	case DefaultChanged:
		return ImpactRisky

	default:
		return ImpactRisky
	}
}
