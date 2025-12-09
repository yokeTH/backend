package model

import "time"

type ClientModel struct {
	ID         int64  `dynamo:"id,hash"`
	SecretHash string `dynamo:"secret_hash"`

	Name        string `dynamo:"name"`
	Description string `dynamo:"description"`

	AllowedCallbackURLs []string `dynamodb:"allowed_callback_urls,set"`

	CreatedAt time.Time `dynamodb:"created_at"`
	CreatedBy int64     `dynamodb:"created_by"`
}
