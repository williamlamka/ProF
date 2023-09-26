package handler

import (
	"errors"
	"new_project/config"
	"new_project/model"
	"new_project/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetBill(userId uuid.UUID, billType string) ([]model.Bill, error) {
	db := config.DB()
	var bill []model.Bill
	result := db.Where("user_id = ?", userId).Where("type = ?", billType).Find(&bill)
	if result.Error != nil {
		return nil, result.Error
	}
	return bill, nil
}

func GetBillById(userId uuid.UUID, id string) (*model.Bill, error) {
	db := config.DB()
	var bill model.Bill
	result := db.Where("user_id = ?", userId).Where("id = ?", id).Find(&bill)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bill, nil
}

func CreateBill(userId uuid.UUID, dto *model.CreateBillDto) (*model.Bill, error) {
	db := config.DB()
	if dto.Type == model.Personal {
		dto.Participant = 1
	}
	currentDate := utils.CurrentTimeWithLocalTZ()
	newBill := model.Bill{
		Type:        dto.Type,
		Description: dto.Description,
		Participant: dto.Participant,
		Price:       utils.LimitDecimalDigits(dto.Price),
		Plan:        dto.Plan,
		CreatedAt:   currentDate,
		ModifiedAt:  currentDate,
		LastPaidAt:  currentDate,
		UserID:      userId,
	}
	result := db.Create(&newBill)
	if result.Error != nil {
		return nil, errors.New(result.Error.Error())
	}
	newTransaction := &model.CreateTransactionDto{
		Category:        	"Bill",
		Description:     	dto.Description,
		Price:           	dto.Price,
		TransactionDate: 	currentDate.Format(config.DateFormat),
	}
	CreateTransaction(userId, newTransaction)
	return &newBill, nil
}

func UpdateBill(userId uuid.UUID, dto *model.UpdateBillDto, id string) error {
	db := config.DB()
	result := db.Where("user_id = ?", userId).Where("id = ?", id).Updates(model.Bill{
		Type:        dto.Type,
		Description: dto.Description,
		Participant: dto.Participant,
		Price:       utils.LimitDecimalDigits(dto.Price),
		Plan:        dto.Plan,
		ModifiedAt:  utils.CurrentTimeWithLocalTZ(),
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func RemoveBill(userId uuid.UUID, c echo.Context) (*model.BillType, error) {
	db := config.DB()
	id := c.Param("id")
	var deletedBill model.Bill
	result := db.Where("id = ?", id).Where("user_id = ?", userId).Find(&deletedBill)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New(config.Error)
		}
		return nil, result.Error
	}
	billType := deletedBill.Type
	result = db.Delete(&deletedBill)
	if result.Error != nil {
		return nil, result.Error
	}
	return &billType, nil
}

func Subscribe(userId uuid.UUID) error {
	db := config.DB()
	var billList []model.Bill
	currentDate := utils.CurrentTimeWithLocalTZ()
	currentYear := currentDate.Format("2006")
	currentMonth := currentDate.Format("01")
	result := db.Where("user_id = ?", userId).Find(&billList)
	if result.Error != nil {
		return result.Error
	}
	for _, element := range billList {
		year := element.LastPaidAt.Format("2006")
		month := element.LastPaidAt.Format("01")
		if month != currentMonth && year <= currentYear {
			newTransaction := &model.CreateTransactionDto{
				Category:        	"Bill",
				Description:     	element.Description,
				Price:           	element.Price / float32(element.Participant),
				TransactionDate: 	currentDate.Format(config.DateFormat),
			}
			CreateTransaction(userId, newTransaction)
			result = db.Where("user_id = ?", userId).Where("id = ?", element.ID).Updates(model.Bill{
				LastPaidAt: currentDate,
			})
			if result.Error != nil {
				return result.Error
			}
		}
	}
	return nil
}
