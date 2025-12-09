package dynamodb

import (
	"context"
	"time"

	"github.com/guregu/dynamo/v2"
	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type accountDynamodb struct {
	table dynamo.Table
}

func NewAccountDynamodb(db *dynamo.DB) *accountDynamodb {
	table := db.Table("accounts")
	dynamo := accountDynamodb{
		table: table,
	}

	return &dynamo
}

func (d *accountDynamodb) Create(ctx context.Context, account *model.AccountModel) error {
	if account.CreatedAt.IsZero() {
		account.CreatedAt = time.Now().UTC()
	}

	return d.table.Put(account).Run(ctx)
}

func (d *accountDynamodb) GetByUserID(ctx context.Context, userID int64) ([]model.AccountModel, error) {
	var accounts []model.AccountModel
	if err := d.table.Scan().Filter("'user_id' = ?", userID).All(ctx, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (d *accountDynamodb) GetByProviderAndProviderUserID(ctx context.Context, provider, providerUserID string) (*model.AccountModel, error) {
	var account model.AccountModel
	err := d.table.
		Get("provider", provider).
		Range("provider_user_id", dynamo.Equal, providerUserID).
		One(ctx, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
