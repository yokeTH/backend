package repository

import (
	"context"

	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type UserPermission interface {
	Create(ctx context.Context, up *model.UserPermissionsModel) error
	GetByUserIDAndClient(ctx context.Context, userID, clientID int64) ([]model.UserPermissionsModel, error)
	Delete(ctx context.Context, userID, clientID int64, name string) error
}
