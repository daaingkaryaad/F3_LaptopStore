package model

type RolePermission struct {
	RoleID       string `json:"role_id" bson:"role_id"`
	PermissionID string `json:"permission_id" bson:"permission_id"`
}
