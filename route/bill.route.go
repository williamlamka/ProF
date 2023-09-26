package route

import (
	"encoding/json"
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

type BillResponse struct {
	ID			uint64			`json:"id"`
	Type        model.BillType	`json:"type"`
	Description string			`json:"description"`
	Participant uint8			`json:"participant"`
	Price       float32			`json:"price"`
	Plan        model.BillPlan	`json:"plan"`
	LastPaidAt  string			`json:"lastPaidAt"`
}

func BillRouteInit(e *echo.Echo) {
	g := e.Group("/bill")
	g.Use(echojwt.WithConfig(config.JwtConfig()))

	g.GET("/all/:type", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		billType := c.Param("type")
		key := fmt.Sprintf("user:%s:bill-%s-type", claims.ID, billType)
		var bills []model.Bill
		val, err := redis.Get(key)
		if err != nil {
			return err
		}
		if val != nil {
			json.Unmarshal(val, &bills)
		} else {
			bills, err = handler.GetBill(claims.ID, billType)
			if err != nil {
				if err.Error() != "record not found" {
					return err
				}
				return c.JSON(http.StatusOK, utils.SuccessResponse{
					Success: true,
				})
			}
			jsonstr, _ := json.Marshal(bills)
			err := redis.Set(key, string(jsonstr))
			if err != nil {
				return err
			}
		}	
		var data []BillResponse
		assignToBillResponseArray(&data, bills)
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
			Data:    data,
		})
	})

	g.GET("/:id", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		id := c.Param("id")
		key := fmt.Sprintf("bill:%s", id)
		var bill *model.Bill
		val, err := redis.Get(key); if err != nil {
			return err
		}
		if val != nil {
			json.Unmarshal(val, &bill)
		} else {
			bill, err = handler.GetBillById(claims.ID, id)
			if err != nil {
				return err
			}
			jsonstr, _ := json.Marshal(bill)
			err = redis.Set(key, string(jsonstr))
			if err != nil {
				return err
			}
		}
		var data BillResponse
		assignToBillResponse(&data, bill)
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
			Data:    data,
		})
	})

	g.POST("", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		var dto model.CreateBillDto
		err := c.Bind(&dto); if err != nil {
			return err
		}
		err = validator.New().Struct(dto)
		if err != nil {
			return err
		}
		if dto.Type == model.Shared {
			dto.Price /= float32(dto.Participant)
		}
		key := fmt.Sprintf("user:%s:bill-%s-type", claims.ID, dto.Type)
		newBill, err := handler.CreateBill(claims.ID, &dto)
		if err != nil {
			if err.Error() != "record not found" {
				fmt.Println(err.Error())
				return err
			}
			return c.JSON(http.StatusOK, utils.SuccessResponse{
				Success: true,
			})
		}
		var data BillResponse
		assignToBillResponse(&data, newBill)
		redis.Delete(key)
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
			Data:    data,
		})
	})

	g.POST("/update/:id", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		id := c.Param("id")
		var dto model.UpdateBillDto
		err := c.Bind(&dto); if err != nil {
			return err
		}
		err = validator.New().Struct(dto)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("user:%s:bill-shared-type", claims.ID)
		key2 := fmt.Sprintf("user:%s:bill-personal-type", claims.ID)
		err = handler.UpdateBill(claims.ID, &dto, id)
		if err != nil {
			return err
		}
		redis.Delete(key)
		redis.Delete(key2)
		return c.JSON(http.StatusOK, utils.SuccessResponse {
			Success: true,
		})
	})

	g.POST("/delete/:id", func(c echo.Context) error {
		claims := config.ExtractJwt(c)
		billType, err := handler.RemoveBill(claims.ID, c)
		key := fmt.Sprintf("user:%s:bill-%s-type", claims.ID, string(*billType))
		if err != nil {
			return err
		}
		redis.Delete(key)
		return c.JSON(http.StatusOK, utils.SuccessResponse{
			Success: true,
		})
	})
}

func assignToBillResponse(data *BillResponse, bill *model.Bill) {
	*data = BillResponse{
		ID: 		 bill.ID,		
		Type:        bill.Type,
		Description: bill.Description,
		Participant: bill.Participant,
		Price:       bill.Price,
		Plan:        bill.Plan,
		LastPaidAt:  bill.LastPaidAt.Format(config.DateFormat),
	}
}

func assignToBillResponseArray(data *[]BillResponse, bills []model.Bill) {
	for _, bill := range bills {
		*data = append(*data, BillResponse{
			ID: 		 bill.ID,		
			Type:        bill.Type,
			Description: bill.Description,
			Participant: bill.Participant,
			Price:       bill.Price,
			Plan:        bill.Plan,
			LastPaidAt:  bill.LastPaidAt.Format(config.DateFormat),
		})
	}
}
