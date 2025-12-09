package repository

import (
	"context"

	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type AccountRepository interface {
	Create(ctx context.Context, account *model.AccountModel) error
	GetByUserID(ctx context.Context, userID int64) ([]model.AccountModel, error)
	GetByProviderAndProviderUserID(ctx context.Context, provider, providerUserID string) (*model.AccountModel, error)
}
