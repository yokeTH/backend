package repository

import (
	"context"

	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type SessionChainRepository interface {
	Create(ctx context.Context, chain *model.SessionChainModel) error
	GetByID(ctx context.Context, chainID string) (*model.SessionChainModel, error)
	Revoke(ctx context.Context, chainID, reason string) error
}
