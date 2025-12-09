package model

type UserPermissionsModel struct {
	UserID         int64  `dynamo:"user_id,hash"`
	ClientID       int64  `dynamo:"client_id"`
	PermissionName string `dynamo:"permission_name"`
}
