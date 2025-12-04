package port

import (
	"context"

	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type JWKGenerator interface {
	Generate(ctx context.Context) (*model.JWKKeyModel, error)
}
