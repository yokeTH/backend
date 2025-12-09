package repository

import (
	"context"

	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type ServiceAccountRepository interface {
	Create(ctx context.Context, sa *model.ServiceAccountModel) error
	GetByID(ctx context.Context, id int64) (*model.ServiceAccountModel, error)
	GetByCreatorID(ctx context.Context, id int64) ([]model.ServiceAccountModel, error)
	Delete(ctx context.Context, id int64) error
}
