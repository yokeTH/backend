package dynamodb

import (
	"context"
	"time"

	"github.com/guregu/dynamo/v2"
	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type clientDynamodb struct {
	table dynamo.Table
}

func NewClientDynamodb(db *dynamo.DB) *clientDynamodb {
	table := db.Table("clients")
	return &clientDynamodb{
		table: table,
	}
}

func (d *clientDynamodb) Create(ctx context.Context, c *model.ClientModel) error {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}
	return d.table.Put(c).Run(ctx)
}

func (d *clientDynamodb) GetByID(ctx context.Context, id int64) (*model.ClientModel, error) {
	var client model.ClientModel
	err := d.table.Get("id", id).One(ctx, &client)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (d *clientDynamodb) GetByCreatorID(ctx context.Context, id int64) ([]model.ClientModel, error) {
	var clients []model.ClientModel
	err := d.table.Scan().Filter("'created_by' = ?", id).All(ctx, &clients)
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func (d *clientDynamodb) Delete(ctx context.Context, id int64) error {
	return d.table.Delete("id", id).Run(ctx)
}
