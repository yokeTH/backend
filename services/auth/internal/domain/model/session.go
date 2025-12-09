package model

import "time"

type SessionModel struct {
	// PK: chain_id (hash) + session_id (range)
	ChainID   string `dynamo:"chain_id,hash"`
	SessionID string `dynamo:"session_id,range"`

	UserID int64 `dynamo:"user_id"`

	// Store only the hash of the refresh token
	RefreshTokenHash string `dynamo:"refresh_token_hash"`

	ParentSessionID     string `dynamo:"parent_session_id,omitempty"`
	ReplacedBySessionID string `dynamo:"replaced_by_session_id,omitempty"`

	IP        string `dynamo:"ip,omitempty"`
	UserAgent string `dynamo:"user_agent,omitempty"`

	ExpiresAt time.Time  `dynamo:"expires_at"`
	RevokedAt *time.Time `dynamo:"revoked_at,omitempty"`

	RevokedReason string    `dynamo:"revoked_reason,omitempty"`
	CreatedAt     time.Time `dynamo:"created_at"`
	UpdatedAt     time.Time `dynamo:"updated_at"`
}
