package drift

import (
	"context"

	"github.com/Bibhu20031/SchemaWatch/internal/snapshot"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Process(
	ctx context.Context,
	schemaID int64,
	fromVersion int,
	toVersion int,
	prev []snapshot.Column,
	curr []snapshot.Column,
) ([]ClassifiedChange, error) {

	changes := Detect(prev, curr)
	classified := Classify(changes)

	if len(classified) == 0 {
		return classified, nil
	}

	err := s.repo.StoreEvents(
		ctx,
		schemaID,
		fromVersion,
		toVersion,
		classified,
	)
	if err != nil {
		return nil, err
	}

	return classified, nil
}
