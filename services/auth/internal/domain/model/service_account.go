package model

import "time"

type ServiceAccountModel struct {
	ID           int64
	Email        string
	PasswordHash string

	CreatedAt time.Time `dynamodb:"created_at"`
	CreatedBy int64     `dynamodb:"created_by"`
}
