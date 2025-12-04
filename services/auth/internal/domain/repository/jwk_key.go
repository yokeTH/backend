package repository

import (
	"context"

	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type JWKKeyRepository interface {
	Count(ctx context.Context) (int, error)
	CreateKey(ctx context.Context, key *model.JWKKeyModel) error
	GetActiveKey(ctx context.Context) (*model.JWKKeyModel, error)
	GetPublicKeys(ctx context.Context) ([]string, error)
	Rotete(ctx context.Context, new *model.JWKKeyModel) error
}
