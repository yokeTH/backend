package dynamodb

import (
	"context"
	"time"

	"github.com/guregu/dynamo/v2"
	"github.com/yokeTH/backend/services/auth/internal/domain/model"
	"github.com/yokeTH/backend/services/auth/internal/domain/repository"
)

type sessionDynamodb struct {
	table dynamo.Table
}

func NewSessionDynamodb(db *dynamo.DB) repository.SessionRepository {
	table := db.Table("sessions")

	return &sessionDynamodb{
		table: table,
	}
}

func (d *sessionDynamodb) Create(ctx context.Context, session *model.SessionModel) error {
	now := time.Now().UTC()
	if session.CreatedAt.IsZero() {
		session.CreatedAt = now
	}
	session.UpdatedAt = now

	return d.table.Put(session).Run(ctx)
}

func (d *sessionDynamodb) Update(ctx context.Context, session *model.SessionModel) error {
	session.UpdatedAt = time.Now().UTC()

	return d.table.Put(session).Run(ctx)
}

func (d *sessionDynamodb) GetByChainAndID(ctx context.Context, chainID, sessionID string) (*model.SessionModel, error) {
	var session model.SessionModel
	err := d.table.
		Get("chain_id", chainID).
		Range("session_id", dynamo.Equal, sessionID).
		One(ctx, &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (d *sessionDynamodb) GetByChainID(ctx context.Context, chainID string) ([]model.SessionModel, error) {
	var sessions []model.SessionModel
	if err := d.table.
		Get("chain_id", chainID).
		All(ctx, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}
