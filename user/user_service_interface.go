package user

import (
	"eventorganizer/golang/models"
)

type UserService interface {
	Register(user *models.User) (*models.User, error)
	GetAllUser() ([]*models.UserNoPassword,error)
	GetUserByID(id int) (*models.UserNoPassword,error)
	DeleteUser(id int) (*models.User, error)
	IsAnyEmailUser(email string) bool
	GetUserByEmail(email string) (*models.User, error)
	UpgradeUser(user *models.User) (*models.UserNoPassword, error)
}
