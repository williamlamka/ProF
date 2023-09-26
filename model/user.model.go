package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID 				uuid.UUID		`gorm:"primaryKey;not null" json:"id"`
	Email 			string 			`gorm:"not null;unique" json:"email"`
	Password		string			`gorm:"not null" json:"password"`
	Username		string			`gorm:"not null" json:"username"`
	Category		pq.StringArray	`gorm:"type:text[]" json:"category"`
	CreatedAt       time.Time 		`gorm:"not null" json:"createdAt"`
	ModifiedAt      time.Time 		`gorm:"not null" json:"modifiedAt"`
}	

type CategoryDto struct {
	Category	string	`json:"category" validate:"required"`
}

type UpdateUserDto struct {
	Username string `json:"username" validate:"required"`
}