package models

import (
	"github.com/google/uuid"
)

type User struct {
	Model
	ID        uuid.UUID `json:"id" gorm:"column:id;primary_key;type:uuid;default:uuid_generate_v4()"`
	FirstName string    `json:"first_name" gorm:"column:first_name"`
	LastName  string    `json:"last_name" gorm:"column:last_name"`
	Email     string    `json:"email" gorm:"column:email;index:idx_name,unique"`
	Password  string    `json:"-" gorm:"column:password"`
	IsActive  bool      `json:"is_active" gorm:"column:is_active;defult=false"`
	Deleted   bool      `json:"deleted" gorm:"column:deleted;defult=false"`
}
