package dynamodb

import (
	"context"
	"time"

	"github.com/guregu/dynamo/v2"
	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type clientPermissionDynamodb struct {
	table dynamo.Table
}

func NewClientPermissionDynamodb(db *dynamo.DB) *clientPermissionDynamodb {
	table := db.Table("client_permissions")
	dynamo := clientPermissionDynamodb{
		table: table,
	}

	return &dynamo
}

func (d *clientPermissionDynamodb) Create(ctx context.Context, cp *model.ClientPermissionModel) error {
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = time.Now().UTC()
	}

	return d.table.Put(cp).Run(ctx)
}

func (d *clientPermissionDynamodb) GetByClientID(ctx context.Context, id int64) ([]model.ClientPermissionModel, error) {
	var cp []model.ClientPermissionModel
	err := d.table.
		Scan().
		Filter("'client_id' = ?", id).
		All(ctx, &cp)
	if err != nil {
		return nil, err
	}

	return cp, nil
}

func (d *clientPermissionDynamodb) Delete(ctx context.Context, clientID int64, name string) error {
	return d.table.Delete("client_id", clientID).Range("name", name).Run(ctx)
}
