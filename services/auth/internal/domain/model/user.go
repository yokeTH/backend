package model

import "time"

type UserModel struct {
	ID           int64  `dynamo:"id,hash"`
	Email        string `dynamo:"email" index:"email-index,hash"`
	PasswordHash string `dynamo:"password_hash"`

	Name  string `dynamo:"name"`
	Image string `dynamo:"image"`

	CreatedAt time.Time `dynamo:"created_at"`
	UpdatedAt time.Time `dynamo:"updated_at"`
}
