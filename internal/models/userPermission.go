package models

import (
	"time"

	"github.com/google/uuid"
)

type UserPermission struct {
	UserID       uuid.UUID  `json:"user_id" gorm:"type:uuid;primaryKey"`
	PermissionID uuid.UUID  `json:"permission_id" gorm:"type:uuid;primaryKey"`
	CreatedAt    time.Time  `json:"created_at"`
	User         User       `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Permission   Permission `json:"-" gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE"`
}
