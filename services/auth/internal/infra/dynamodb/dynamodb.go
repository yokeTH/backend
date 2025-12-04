package dynamodb

import (
	"context"
	"errors"
	"time"

	"github.com/guregu/dynamo/v2"
	"github.com/yokeTH/backend/services/auth/internal/domain/model"
)

type jwkKeyDynamodb struct {
	table dynamo.Table
}

func NewJWKKeyDynamodb(db *dynamo.DB) *jwkKeyDynamodb {
	table := db.Table("jwk_keys")
	dynamo := jwkKeyDynamodb{
		table: table,
	}

	return &dynamo
}

func (d *jwkKeyDynamodb) Count(ctx context.Context) (int, error) {
	return d.table.Scan().Count(ctx)
}

func (d *jwkKeyDynamodb) CreateKey(ctx context.Context, key *model.JWKKeyModel) error {
	if key.CreatedAt.IsZero() {
		key.CreatedAt = time.Now().UTC()
	}
	return d.table.Put(key).Run(ctx)
}

func (d *jwkKeyDynamodb) GetActiveKey(ctx context.Context) (*model.JWKKeyModel, error) {
	var result model.JWKKeyModel

	if err := d.table.Get("status", model.JWKKeyStatusActive).One(ctx, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *jwkKeyDynamodb) GetPublicKeys(ctx context.Context) ([]string, error) {
	var results []model.JWKKeyModel

	if err := d.table.Scan().Filter("'status' = ? OR 'status' = ?", model.JWKKeyStatusActive, model.JWKKeyStatusRetiring).All(ctx, &results); err != nil {
		return nil, err
	}

	pubKeys := make([]string, len(results))
	for i := range len(results) {
		pubKeys[i] = results[i].PublicJWK
	}

	return pubKeys, nil
}

func (r *jwkKeyDynamodb) Rotete(ctx context.Context, newKey *model.JWKKeyModel) error {
	const (
		retiringGrace = 7 * 24 * time.Hour
		retiredKeep   = 30 * 24 * time.Hour
	)

	now := time.Now().UTC()
	var actives []model.JWKKeyModel
	if err := r.table.
		Scan().
		Filter("'status' = ?", model.JWKKeyStatusActive).
		All(ctx, &actives); err != nil && !errors.Is(err, dynamo.ErrNotFound) {
		return err
	}

	for _, a := range actives {
		notAfter := now.Add(retiringGrace)
		notAfterEpoch := notAfter.Unix()

		if err := r.table.
			Update("kid", a.KID).
			Set("status", model.JWKKeyStatusRetiring).
			Set("rotated_at", now).
			Set("not_after", notAfter).
			Set("not_after_epoch", notAfterEpoch).
			Run(ctx); err != nil {
			return err
		}
	}

	var retiring []model.JWKKeyModel
	err := r.table.
		Scan().
		Filter("'status' = ?", model.JWKKeyStatusRetiring).
		All(ctx, &retiring)

	if err != nil && !errors.Is(err, dynamo.ErrNotFound) {
		return err
	}

	for _, k := range retiring {
		if k.NotAfter == nil {
			continue
		}
		if now.Before(*k.NotAfter) {
			continue
		}

		ttlTime := k.NotAfter.Add(retiredKeep)
		ttlEpoch := ttlTime.Unix()

		if uErr := r.table.
			Update("kid", k.KID).
			Set("status", model.JWKKeyStatusRetired).
			Set("not_after_epoch", ttlEpoch).
			Run(ctx); uErr != nil {
			return uErr
		}
	}

	newKey.Status = model.JWKKeyStatusActive
	if newKey.CreatedAt.IsZero() {
		newKey.CreatedAt = now
	}

	if newKey.NotBefore == nil {
		nb := now
		newKey.NotBefore = &nb
	}
	newKey.NotAfterEpoch = nil

	return r.table.Put(newKey).Run(ctx)
}
