package model

import (
	"time"

	"github.com/google/uuid"
)

type BillType string

const (
    Personal BillType = "personal"
    Shared BillType = "shared"
)

type BillPlan string

const (
    Year BillPlan = "year"
    Season BillPlan = "season"
	Month BillPlan = "month"
)

type Bill struct {
	ID			uint64		`gorm:"primaryKey;autoIncrement" json:"id"`
	Type 		BillType	`gorm:"not null" json:"type"`
	Description	string		`gorm:"not null;unique" json:"description"`
	Participant uint8		`gorm:"not null" json:"participant"`
	Price		float32		`gorm:"not null" json:"price"`
	Plan		BillPlan	`gorm:"not null" json:"plan"`
	CreatedAt   time.Time   `gorm:"not null" json:"createdAt"`
    ModifiedAt  time.Time   `gorm:"not null" json:"modifiedAt"`
	LastPaidAt  time.Time   `gorm:"not null" json:"lastPaidAt"`
    UserID     	uuid.UUID   `gorm:"index not null" json:"userId"`
}

type CreateBillDto struct {
	Type 		BillType `json:"type" validate:"required"`
	Description	string   `json:"description" validate:"required"`
	Participant uint8    `json:"participant" validate:"required"`
	Price		float32	 `json:"price" validate:"required,gt=0"`		 
	Plan		BillPlan `json:"plan" validate:"required"`  
}

type UpdateBillDto struct {
	Type 		BillType `json:"type" validate:"required"`
	Description	string   `json:"description" validate:"required"`
	Participant uint8    `json:"participant" validate:"required"`
	Price		float32	 `json:"price" validate:"required,gt=0"`		 
	Plan		BillPlan `json:"plan" validate:"required"`  
}
