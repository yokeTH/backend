package dynamodb

import (
	"context"
	"time"

	"github.com/guregu/dynamo/v2"
	"github.com/yokeTH/backend/services/auth/internal/domain/model"
	"github.com/yokeTH/backend/services/auth/internal/domain/repository"
)

type sessionChainDynamodb struct {
	table dynamo.Table
}

func NewSessionChainDynamodb(db *dynamo.DB) repository.SessionChainRepository {
	table := db.Table("session_chains")

	return &sessionChainDynamodb{
		table: table,
	}
}

func (d *sessionChainDynamodb) Create(ctx context.Context, chain *model.SessionChainModel) error {
	now := time.Now().UTC()
	if chain.CreatedAt.IsZero() {
		chain.CreatedAt = now
	}
	chain.UpdatedAt = now

	if chain.Status == "" {
		chain.Status = model.SessionChainStatusActive
	}

	return d.table.Put(chain).Run(ctx)
}

func (d *sessionChainDynamodb) GetByID(ctx context.Context, chainID string) (*model.SessionChainModel, error) {
	var chain model.SessionChainModel
	if err := d.table.
		Get("chain_id", chainID).
		One(ctx, &chain); err != nil {
		return nil, err
	}

	return &chain, nil
}

func (d *sessionChainDynamodb) Revoke(ctx context.Context, chainID, reason string) error {
	now := time.Now().UTC()

	return d.table.
		Update("chain_id", chainID).
		Set("status", model.SessionChainStatusRevoked).
		Set("revoked_at", now).
		Set("revoked_reason", reason).
		Set("updated_at", now).
		Run(ctx)
}
