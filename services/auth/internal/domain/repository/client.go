package repository

import (
	"context"

	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type ClientRepository interface {
	Create(ctx context.Context, c *model.ClientModel) error
	GetByID(ctx context.Context, id int64) (*model.ClientModel, error)
	GetByCreatorID(ctx context.Context, id int64) ([]model.ClientModel, error)
	Delete(ctx context.Context, id int64) error
}
