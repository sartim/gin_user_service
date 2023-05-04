package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID          uuid.UUID `json:"id" gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	Name        string    `json:"name" gorm:"index:idx_name,unique"`
	Description string    `json:"description"`
	Deleted     *bool     `json:"deleted" gorm:"default=false"`
}
