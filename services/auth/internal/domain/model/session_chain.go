package model

import "time"

const (
	SessionChainStatusActive  = "ACTIVE"
	SessionChainStatusRevoked = "REVOKED"
)

type SessionChainModel struct {
	// Partition key
	ChainID string `dynamo:"chain_id,hash"`

	UserID int64  `dynamo:"user_id"`
	Status string `dynamo:"status"`

	RevokedAt     *time.Time `dynamo:"revoked_at,omitempty"`
	RevokedReason string     `dynamo:"revoked_reason,omitempty"`

	CreatedAt time.Time `dynamo:"created_at"`
	UpdatedAt time.Time `dynamo:"updated_at"`
}
