package model

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID              	uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Category        	string    `gorm:"not null" json:"category"`
	Description     	string    `gorm:"not null" json:"description"`
	Price           	float32   `gorm:"not null" json:"price"`
	TransactionDate 	string 	  `gorm:"not null" json:"transactionDate"`
	TransactionMonth 	string 	  `gorm:"not null" json:"transactionMonth"`
	TransactionYear 	string 	  `gorm:"not null" json:"transactionYear"`
	CreatedAt       	time.Time `gorm:"not null" json:"createdAt"`
	ModifiedAt      	time.Time `gorm:"not null" json:"modifiedAt"`
	UserID          	uuid.UUID `gorm:"not null" json:"userId"`
}

type TransactionChartData struct {
	Category	string	`json:"category"`
	TotalPrice	float32	`json:"totalPrice"`
}

type CreateTransactionDto struct {
	Category        	string  `json:"category" validate:"required"`
	Description     	string  `json:"description" validate:"required"`
	Price           	float32 `json:"price" validate:"required,gt=0"`
	TransactionDate 	string  `json:"transactionDate" validate:"required"`
}

type UpdateTransactionDto struct {
	Category        	string  `json:"category" validate:"required"`
	Description     	string  `json:"description" validate:"required"`
	Price           	float32 `json:"price" validate:"required,gt=0"`
	TransactionDate 	string  `json:"transactionDate" validate:"required"`
}
