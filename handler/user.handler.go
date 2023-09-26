package handler

import (
	"errors"
	"new_project/config"
	"new_project/model"
	"new_project/utils"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func GetUser(userId uuid.UUID) (*model.User, error) {
	db := config.DB()
	var user model.User
	db.Find(&user, userId)
	return &user, nil
}

func GetUserByEmail(email string) (*model.User, error) {
	db := config.DB()
	var user model.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, errors.New("user does not exist")
	}
	return &user, nil
}

func CreateUser(dto UserRegisterDto) (*model.User, error) {
	db := config.DB()
	newUser := model.User{
		ID:         uuid.New(),
		Email:      dto.Email,
		Password:   dto.Password,
		Username:   dto.UserName,
		Category:   []string{},
		CreatedAt:  utils.CurrentTimeWithLocalTZ(),
		ModifiedAt: utils.CurrentTimeWithLocalTZ(),
	}
	db.Create(&newUser)
	return &newUser, nil
}

func UpdateUser(userId uuid.UUID, dto model.UpdateUserDto) error {
	db := config.DB()
	user := &model.User{ID: userId}
	result := db.Model(user).Updates(&model.User{
		Username:   dto.Username,
		ModifiedAt: utils.CurrentTimeWithLocalTZ(),
	})
	if result.Error != nil {
		return errors.New("user Update Error")
	}
	return nil
}

func GetCategory(userId uuid.UUID) *pq.StringArray {
	db := config.DB()
	var user model.User
	db.Find(&user, userId)
	return &user.Category
}

func AddCategory(userId uuid.UUID, category string) {
	db := config.DB()
	var user model.User
	db.Find(&user, userId)
	if !utils.Contain(user.Category, category) {
		newCategory := append(user.Category, category)
		db.Model(&user).Updates(&model.User{
			Category:   newCategory,
			ModifiedAt: utils.CurrentTimeWithLocalTZ(),
		})
	}
}

func RemoveCategroy(userId uuid.UUID, category string) {
	db := config.DB()
	var user model.User
	db.Find(&user, userId)
	index := utils.Index(user.Category, category)
	if index != -1 {
		newCategory := utils.Remove(user.Category, index)
		db.Model(&user).Updates(&model.User{
			Category:   newCategory,
			ModifiedAt: utils.CurrentTimeWithLocalTZ(),
		})
	}
}
