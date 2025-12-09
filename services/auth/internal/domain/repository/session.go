package repository

import (
	"context"

	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type SessionRepository interface {
	Create(ctx context.Context, session *model.SessionModel) error
	Update(ctx context.Context, session *model.SessionModel) error

	GetByChainAndID(ctx context.Context, chainID, sessionID string) (*model.SessionModel, error)
	GetByChainID(ctx context.Context, chainID string) ([]model.SessionModel, error)
}
