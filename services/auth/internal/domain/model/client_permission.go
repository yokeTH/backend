package model

import "time"

type ClientPermissionModel struct {
	ClientID int64  `dynamo:"client_id,hash"`
	Name     string `dynamo:"name,range"`

	Description string `dynamo:"description"`

	CreatedAt time.Time `dynamodb:"created_at"`
	CreatedBy int64     `dynamodb:"created_by"`
}
