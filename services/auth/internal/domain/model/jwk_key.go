package model

import "time"

type JWKKeyStatus string

const (
	JWKKeyStatusActive   JWKKeyStatus = "ACTIVE"
	JWKKeyStatusRetiring JWKKeyStatus = "RETIRING"
	JWKKeyStatusRetired  JWKKeyStatus = "RETIRED"
)

type JWKKeyModel struct {
	KID string `dynamo:"kid,hash"`

	Alg    string       `dynamo:"alg"`
	Status JWKKeyStatus `dynamo:"status"`

	PublicJWK string `dynamo:"public_jwk"`

	PrivCiphertext []byte `dynamo:"priv_ciphertext"`
	PrivNonce      []byte `dynamo:"priv_nonce"`
	WrappedDEK     []byte `dynamo:"wrapped_dek"`

	KEKRef string `dynamo:"kek_ref,omitempty"`

	CreatedAt time.Time  `dynamo:"created_at"`
	RotatedAt *time.Time `dynamo:"rotated_at,omitempty"`
	NotBefore *time.Time `dynamo:"not_before,omitempty"`
	NotAfter  *time.Time `dynamo:"not_after,omitempty"`

	NotAfterEpoch *int64 `dynamo:"not_after_epoch,omitempty"` // used for TTL
}
