package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ValidateAPIKey(ctx context.Context, rawKey string) (bool, error) {
	hash := hashKey(rawKey)
	return s.repo.IsValidKey(ctx, hash)
}

func hashKey(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}
