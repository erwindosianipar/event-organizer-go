package service

import (
	"eventorganizer/golang/models"
	"eventorganizer/golang/user"
)

type UserService struct {
	UserRepo user.UserRepo
}

func CreateUserService(userRepo user.UserRepo) user.UserService {
	return &UserService{userRepo}
}

func (h *UserService) Register(user *models.User) (*models.UserNoPassword, error) {
	return h.UserRepo.Register(user)
}

func (h *UserService) GetAllUser() ([]*models.UserNoPassword, error) {
	return h.UserRepo.GetAllUser()
}

func (h *UserService) GetUserByID(id int) (*models.UserNoPassword, error) {
	return h.UserRepo.GetUserByID(id)
}

func (h *UserService) DeleteUser(id int) (*models.User, error) {
	return h.UserRepo.DeleteUser(id)
}

func (h *UserService) IsAnyEmailUser(email string) bool {
	return h.UserRepo.IsAnyEmailUser(email)
}

func (h *UserService) GetUserByEmail(email string) (*models.User, error) {
	return h.UserRepo.GetUserByEmail(email)
}

func (h *UserService) UpgradeUser(user *models.User) (*models.UserNoPassword, error) {
	return h.UserRepo.UpgradeUser(user)
}

func (h *UserService) HandleUpgrade(id int, status string) (*models.UserNoPassword, error) {
	return h.UserRepo.HandleUpgrade(id, status)
}


