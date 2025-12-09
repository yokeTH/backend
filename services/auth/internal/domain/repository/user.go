package repository

import (
	"context"

	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type UserRepository interface {
	Create(ctx context.Context, usr *model.UserModel) error
	GetByID(ctx context.Context, id int64) (*model.UserModel, error)
	GetByEmail(ctx context.Context, email string) (*model.UserModel, error)
}
