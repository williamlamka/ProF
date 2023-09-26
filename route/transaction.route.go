package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"new_project/config"
	"new_project/handler"
	"new_project/model"
	"new_project/redis"
	"new_project/utils"

	"github.com/go-playground/validator/v10"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type TransactionResponse struct {
	ID              uint64  `json:"id"`
	Category        string  `json:"category"`
	Description     string  `json:"description"`
	Price           float32 `json:"price"`
	TransactionDate string  `json:"transactionDate"`
	CreatedAt       string  `json:"createdAt"`
	ModifiedAt      string  `json:"modifiedAt"`
}

func TransactionRouteInit(e *echo.Echo) {
	g := e.Group("/transaction")
	g.Use(echojwt.WithConfig(config.JwtConfig()))

	g.GET("/all/:transactionDate", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		transactionDate := c.Param("transactionDate")
		key := fmt.Sprintf("user:%s:transaction-%s", claims.ID, transactionDate)
		var transactions []model.Transaction
		val, err := redis.Get(key); if err != nil {
			return err
		}
		if val != nil {
			json.Unmarshal(val, &transactions)
		} else {
			transactions, err = handler.GetTransaction(claims.ID, transactionDate)
			if err != nil {
				if err.Error() != "record not found" {
					return err
				}
				return c.JSON(http.StatusOK, utils.SuccessResponse{
					Success: true,
				})
			}
			jsonstr, _ := json.Marshal(transactions)
			err := redis.Set(key, string(jsonstr))
			if err != nil {
				return err
			}
		}
		var data []TransactionResponse
		assignToTransactionResponseArray(&data, transactions)
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
			Data:    data,
		})
	})

	g.GET("/latestthree", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		key := fmt.Sprintf("user:%s:transaction-latest-three", claims.ID)
		var transactions []model.Transaction
		val, err := redis.Get(key)
		if err != nil {
			return err
		}
		if val != nil {
			json.Unmarshal(val, &transactions)
		} else {
			transactions, err = handler.GetLatestThreeTransaction(claims.ID, c)
			if err != nil {
				if err.Error() != "record not found" {
					return err
				}
				return c.JSON(http.StatusOK, utils.SuccessResponse{
					Success: true,
				})
			}
			jsonstr, _ := json.Marshal(transactions)
			err := redis.Set(key, string(jsonstr))
			if err != nil {
				return err
			}
		}
		var data []TransactionResponse
		assignToTransactionResponseArray(&data, transactions)
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
			Data:    data,
		})
	})

	g.GET("/chartData/:date", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		requiredDate := c.Param("date")
		transactionChartData, err := handler.GetTransactionChartData(claims.ID, requiredDate)
		if err != nil {
			if err.Error() != "record not found" {
				return err
			}
			return c.JSON(http.StatusOK, utils.SuccessResponse{
				Success: true,
			})
		}
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
			Data:    transactionChartData,
		})
	})

	g.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")
		key := fmt.Sprintf("transaction:%s", id)
		val, err := redis.Get(key)
		if err != nil {
			return err
		}
		var transaction *model.Transaction
		if val != nil {
			json.Unmarshal(val, &transaction)
		} else {
			transaction, err = handler.GetTransactionById(c, id)
			if err != nil {
				if err.Error() != "record not found" {
					return err
				}
				return c.JSON(http.StatusOK, utils.SuccessResponse{
					Success: true,
				})
			}
			jsonstr, _ := json.Marshal(transaction)
			err = redis.Set(key, string(jsonstr))
			if err != nil {
				return err
			}
		}
		var data TransactionResponse
		assignToTransactionResponse(&data, transaction)
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
			Data:    data,
		})
	})

	g.POST("", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		key := fmt.Sprintf("user:%s:transaction-latest-three", claims.ID)
		var dto model.CreateTransactionDto
		err := c.Bind(&dto)
		if err != nil {
			return err
		}
		err = validator.New().Struct(dto)
		if err != nil {
			return err
		}
		date := utils.GetTransactionDate(dto.TransactionDate)
		key2 := fmt.Sprintf("user:%s:transaction-%s", claims.ID, date)
		newTransaction, err := handler.CreateTransaction(claims.ID, &dto)
		if err != nil {
			return errors.New(err.Error())
		}
		var data TransactionResponse
		assignToTransactionResponse(&data, newTransaction)
		redis.Delete(key)
		redis.Delete(key2)
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
			Data:    data,
		})
	})

	g.POST("/update/:id", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		id := c.Param("id")
		key := fmt.Sprintf("transaction:%s", id)
		key2 := fmt.Sprintf("user:%s:transaction-latest-three", claims.ID)
		var dto model.UpdateTransactionDto
		err := c.Bind(&dto)
		if err != nil {
			return err
		}
		err = validator.New().Struct(dto)
		if err != nil {
			return err
		}
		date := utils.GetTransactionDate(dto.TransactionDate)
		key3 := fmt.Sprintf("user:%s:transaction-%s", claims.ID, date)
		err = handler.UpdateTransaction(claims.ID, &dto, id)
		if err != nil {
			if err.Error() == "record not found" {
				return errors.New(config.Error)
			}
			return err
		}
		redis.Delete(key)
		redis.Delete(key2)
		redis.Delete(key3)
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
		})
	})

	g.POST("/delete/:id", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		id := c.Param("id")
		key := fmt.Sprintf("transaction:%s", id)
		key2 := fmt.Sprintf("user:%s:transaction-latest-three", claims.ID)
		date, err := handler.DeleteTransaction(claims.ID, id)
		if err != nil {
			if err.Error() == "record not found" {
				return errors.New(config.Error)
			}
			return err
		}
		key3 := fmt.Sprintf("user:%s:transaction-%s", claims.ID, *date)
		redis.Delete(key)
		redis.Delete(key2)
		redis.Delete(key3)
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
		})
	})
}

func assignToTransactionResponse(data *TransactionResponse, transaction *model.Transaction) {
	*data = TransactionResponse{
		ID:              transaction.ID,
		Category:        transaction.Category,
		Description:     transaction.Description,
		Price:           transaction.Price,
		TransactionDate: transaction.TransactionDate,
		CreatedAt:       transaction.CreatedAt.Format(config.DateTimeFormat),
		ModifiedAt:      transaction.ModifiedAt.Format("2006-01-02 15:04:05"),
	}
}

func assignToTransactionResponseArray(data *[]TransactionResponse, transactions []model.Transaction) {
	for _, transaction := range transactions {
		*data = append(*data, TransactionResponse{
			ID:              transaction.ID,
			Category:        transaction.Category,
			Description:     transaction.Description,
			Price:           transaction.Price,
			TransactionDate: transaction.TransactionDate,
			CreatedAt:       transaction.CreatedAt.Format(config.DateTimeFormat),
			ModifiedAt:      transaction.ModifiedAt.Format("2006-01-02 15:04:05"),
		})
	}
}
