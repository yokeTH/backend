package model

import "time"

type AccountModel struct {
	UserID   int64  `dynamo:"user_id,hash"`
	Provider string `dynamo:"provider,range" index:"provider-user-index,hash"`

	ProviderUserID string `dynamo:"provider_user_id" index:"provider-user-index,range"`

	AccessToken           string    `dynamo:"access_token"`
	AccessTokenExpiresAt  time.Time `dynamo:"access_token_expires_at"`
	RefreshToken          string    `dynamo:"refresh_token"`
	RefreshTokenExpiresAt time.Time `dynamo:"refresh_token_expires_at"`
	IDToken               string    `dynamo:"id_token"`
	Scopes                string    `dynamo:"scopes"`

	CreatedAt time.Time `dynamo:"created_at"`
	UpdatedAt time.Time `dynamo:"updated_at"`
}
