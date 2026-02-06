package drift

import (
	"context"

	"github.com/Bibhu20031/SchemaWatch/internal/notify"
	"github.com/Bibhu20031/SchemaWatch/internal/schema"
	"github.com/Bibhu20031/SchemaWatch/internal/snapshot"
)

type Service struct {
	repo       *Repository
	schemaRepo *schema.Repository
}

func NewService(
	repo *Repository,
	schemaRepo *schema.Repository,
) *Service {
	return &Service{repo: repo, schemaRepo: schemaRepo}
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

	if err := s.repo.StoreEvents(
		ctx,
		schemaID,
		fromVersion,
		toVersion,
		classified,
	); err != nil {
		return nil, err
	}

	for _, c := range classified {
		if c.Impact != ImpactBreaking {
			continue
		}

		url, err := s.schemaRepo.GetWebhookURL(ctx, schemaID)
		if err != nil || url == "" {
			continue
		}

		_ = notify.Send(url, notify.Payload{
			SchemaID:    schemaID,
			Impact:      string(c.Impact),
			Summary:     string(c.Type) + " on " + c.ColumnName,
			VersionFrom: fromVersion,
			VersionTo:   toVersion,
		})
	}

	return classified, nil
}
