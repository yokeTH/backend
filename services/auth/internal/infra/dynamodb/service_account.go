package dynamodb

import (
	"context"
	"time"

	"github.com/guregu/dynamo/v2"
	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type serviceAccountDynamodb struct {
	table dynamo.Table
}

func NewServiceAccountDynamodb(db *dynamo.DB) *serviceAccountDynamodb {
	table := db.Table("service_accounts")
	return &serviceAccountDynamodb{
		table: table,
	}
}

func (d *serviceAccountDynamodb) Create(ctx context.Context, sa *model.ServiceAccountModel) error {
	if sa.CreatedAt.IsZero() {
		sa.CreatedAt = time.Now().UTC()
	}
	return d.table.Put(sa).Run(ctx)
}

func (d *serviceAccountDynamodb) GetByID(ctx context.Context, id int64) (*model.ServiceAccountModel, error) {
	var sa model.ServiceAccountModel
	err := d.table.Get("id", id).One(ctx, &sa)
	if err != nil {
		return nil, err
	}
	return &sa, nil
}

func (d *serviceAccountDynamodb) GetByCreatorID(ctx context.Context, id int64) ([]model.ServiceAccountModel, error) {
	var sas []model.ServiceAccountModel
	err := d.table.Scan().Filter("'created_by' = ?", id).All(ctx, &sas)
	if err != nil {
		return nil, err
	}
	return sas, nil
}

func (d *serviceAccountDynamodb) Delete(ctx context.Context, id int64) error {
	return d.table.Delete("id", id).Run(ctx)
}
