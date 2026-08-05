package models

import (
	"time"

	"github.com/google/uuid"
)

type RolePermission struct {
	RoleID       uuid.UUID  `json:"role_id" gorm:"type:uuid;primaryKey"`
	PermissionID uuid.UUID  `json:"permission_id" gorm:"type:uuid;primaryKey"`
	CreatedAt    time.Time  `json:"created_at"`
	Role         Role       `json:"-" gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`
	Permission   Permission `json:"-" gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE"`
}
