package repo

import (
	"errors"
	"eventorganizer/golang/models"
	"eventorganizer/golang/user"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var _ user.UserRepo = &UserRepoPostgreImpl{}

type UserRepoPostgreImpl struct {
	db *gorm.DB
}

func CreateUserRepoPostgreImpl(db *gorm.DB) user.UserRepo {
	return &UserRepoPostgreImpl{db}
}

func (h *UserRepoPostgreImpl) Register(user *models.User) (*models.UserNoPassword, error) {
	if err := h.db.Table("users").Save(&user).Error; err != nil {
		logrus.Error(err)
		return nil, errors.New("ERROR: Error when insert data user into database")
	}

	return &models.UserNoPassword{
		OrmModel: models.OrmModel{
			ID: user.OrmModel.ID,
		},
		Email:            user.Email,
		Name:             user.Name,
		Avatar:           user.Avatar,
		Role:             user.Role,
		SubmissionStatus: user.SubmissionStatus,
		EventOrganizer:   models.EventOrganizer{},
	}, nil
}

func (h *UserRepoPostgreImpl) GetAllUser() ([]*models.UserNoPassword, error) {
	userList := make([]*models.UserNoPassword, 0)

	if err := h.db.Table("users").Find(&userList).Error; err != nil {
		logrus.Error(err)
		return nil, errors.New("ERROR: Error when get all data users")
	}

	return userList, nil
}

func (h *UserRepoPostgreImpl) GetUserByID(id int) (*models.UserNoPassword, error) {
	dataUser := new(models.UserNoPassword)

	if err := h.db.Table("users").Where("id = ?", id).First(&dataUser).Error; err != nil {
		logrus.Error(err)
		return nil, errors.New("ERROR: Error no data user with id you entered")
	}

	return dataUser, nil
}

func (h *UserRepoPostgreImpl) DeleteUser(id int) (*models.User, error) {
	if err := h.db.Table("users").Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		logrus.Error(err)
		return nil, errors.New("ERROR: Error when delete data user")
	}

	return nil, nil
}

func (h *UserRepoPostgreImpl) IsAnyEmailUser(email string) bool {
	var total int
	h.db.Table("users").Where("email = ?", email).Count(&total)

	if total > 0 {
		return true
	}

	return false
}

func (h *UserRepoPostgreImpl) GetUserByEmail(email string) (*models.User, error) {
	dataUser := new(models.User)

	if err := h.db.Table("users").Where("email = ?", email).First(&dataUser).Error; err != nil {
		logrus.Error(err)
		return nil, errors.New("ERROR: Error when get data user by email")
	}

	return dataUser, nil
}

func (h *UserRepoPostgreImpl) UpgradeUser(user *models.User) (*models.UserNoPassword, error) {
	h.db.Table("users").Where("id = ?", user.ID).Updates(map[string]interface{}{
		"submission_status": "submitted",
		"name_eo":           user.NameEo,
		"ktp_number":        user.KTPNumber,
		"ktp_photo":         user.KTPPhoto,
		"siup_number":       user.SIUPNumber,
	})

	dataUser, err := h.GetUserByID(int(user.ID))
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("ERROR: Error when get data user by id")
	}

	return dataUser, nil
}

func (h *UserRepoPostgreImpl) HandleUpgrade(id int, status string) (*models.UserNoPassword, error) {
	isVerify := false
	role := "user"
	if status == "accepted" {
		role = "event_organizer"
		isVerify = true
	}

	dataUser, _ := h.GetUserByID(id)
	if dataUser.SubmissionStatus == "not_submit" {
		return nil, errors.New("ERROR: This user is not submit any request yet")
	}

	h.db.Table("users").Where("id = ?", id).Updates(map[string]interface{}{
		"role":              role,
		"submission_status": status,
		"is_verify":         isVerify,
	})

	return dataUser, nil
}


