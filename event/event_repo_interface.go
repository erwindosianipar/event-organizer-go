package event

import (
	"eventorganizer/golang/models"
	"github.com/jinzhu/gorm"
)

type EventRepo interface {
	BeginTrans()*gorm.DB
	AddEvent (event *models.Event,tx *gorm.DB)error
	GetAllEvent()(*[]models.Event,error)
}