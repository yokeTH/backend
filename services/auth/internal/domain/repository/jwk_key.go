package repository

import "context"

type JWKKeyRepository interface {
	Count(ctx context.Context) (int, error)
	CreateKey(ctx context.Context, key JWKKeyRepository) error
	GetActiveKey(ctx context.Context) (*JWKKeyRepository, error)
	GetPublicKeys(ctx context.Context) ([]JWKKeyRepository, error)
	Rotete(ctx context.Context) error
}
