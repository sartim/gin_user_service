package models

import (
	"gorm.io/gorm"
)

type UserPermission struct {
	gorm.Model
	Users       []User
	Permissions []Permission
}
