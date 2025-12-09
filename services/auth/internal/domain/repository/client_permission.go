package repository

import (
	"context"

	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type ClientPermissionRepository interface {
	Create(ctx context.Context, cp *model.ClientPermissionModel) error
	GetByClientID(ctx context.Context, id int64) ([]model.ClientPermissionModel, error)
	Delete(ctx context.Context, clientID int64, name string) error
}
