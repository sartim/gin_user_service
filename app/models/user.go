package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        uint   `json:"id" gorm:"primary_key"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsActive  bool   `json:"is_active"`
	Deleted   bool   `json:"deleted"`
}
