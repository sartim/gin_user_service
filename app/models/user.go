package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email" gorm:"index:idx_name,unique"`
	Password  string    `json:"password"`
	IsActive  bool      `json:"is_active" gorm:"defult=false"`
	Deleted   bool      `json:"deleted" gorm:"defult=false"`
}
