package handler

import (
	"errors"
	"fmt"
	"new_project/config"
	"new_project/model"
	"new_project/utils"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetTransaction(userId uuid.UUID, date string) ([]model.Transaction, error) {
	db := config.DB()
	date = fmt.Sprintf("%%%s%%", date)
	var transaction []model.Transaction
	result := db.Where("user_id = ?", userId).
		Where("transaction_date LIKE ?", date).
		Order("transaction_date desc, created_at desc").
		Find(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return transaction, nil
}

func GetLatestThreeTransaction(userId uuid.UUID, c echo.Context) ([]model.Transaction, error) {
	db := config.DB()
	var transaction []model.Transaction
	result := db.Where("user_id = ?", userId).
		Order("transaction_date desc, created_at desc").
		Limit(3).
		Find(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return transaction, nil
}

func GetTransactionById(c echo.Context, id string) (*model.Transaction, error) {
	db := config.DB()
	var transaction model.Transaction
	result := db.Where("id = ?", id).First(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return &transaction, nil
}

func GetTransactionChartData(userId uuid.UUID, requiredDate string) ([]model.TransactionChartData, error) {
	db := config.DB()
	requiredYear := strings.Split(requiredDate, "-")[0]
	requiredMonth := strings.Split(requiredDate, "-")[1]
	var transactionChartData []model.TransactionChartData
	result := db.Model(&model.Transaction{}).
		Select("category, sum(price) as total_price").
		Where("user_Id = ?", userId).
		Where("transaction_year", requiredYear).
		Where("transaction_month", requiredMonth).
		Group("category").
		Find(&transactionChartData)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactionChartData, nil
}

func CreateTransaction(userId uuid.UUID, dto *model.CreateTransactionDto) (*model.Transaction, error) {
	db := config.DB()
	AddCategory(userId, dto.Category)
	newTransaction := model.Transaction{
		Category:         dto.Category,
		Description:      dto.Description,
		Price:            utils.LimitDecimalDigits(dto.Price),
		TransactionDate:  dto.TransactionDate,
		TransactionYear:  strings.Split(dto.TransactionDate, "-")[0],
		TransactionMonth: strings.Split(dto.TransactionDate, "-")[1],
		CreatedAt:        utils.CurrentTimeWithLocalTZ(),
		ModifiedAt:       utils.CurrentTimeWithLocalTZ(),
		UserID:           userId,
	}
	result := db.Create(&newTransaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newTransaction, nil
}

func UpdateTransaction(userId uuid.UUID, dto *model.UpdateTransactionDto, id string) error {
	db := config.DB()
	result := db.Where("user_id = ?", userId).Where("id = ?", id).Updates(model.Transaction{
		Category:         dto.Category,
		Description:      dto.Description,
		Price:            utils.LimitDecimalDigits(dto.Price),
		TransactionDate:  dto.TransactionDate,
		TransactionYear:  strings.Split(dto.TransactionDate, "-")[0],
		TransactionMonth: strings.Split(dto.TransactionDate, "-")[1],
		ModifiedAt:       utils.CurrentTimeWithLocalTZ(),
	})
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return errors.New(config.Error)
		}
		return result.Error
	}
	return nil
}

func DeleteTransaction(userId uuid.UUID, id string) (*string, error) {
	db := config.DB()
	var transaction model.Transaction
	result := db.Where("id = ?", id).Where("user_id = ?", userId).First(&transaction)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New(config.Error)
		}
		return nil, result.Error
	}
	date := utils.GetTransactionDate(transaction.TransactionDate)
	result = db.Delete(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return &date, nil
}
