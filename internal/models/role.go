package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `json:"name" gorm:"index:idx_name,unique"`
	Description string `json:"description"`
	Deleted     *bool  `json:"deleted" gorm:"default=false"`
}
