package models

import (
	"gorm.io/gorm"
)

type RolePermission struct {
	gorm.Model
	Roles       []Role
	Permissions []Permission
}
