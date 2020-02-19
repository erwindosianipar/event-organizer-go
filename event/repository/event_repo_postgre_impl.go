package repository

import (
	"eventorganizer/golang/event"
	"eventorganizer/golang/models"
	"github.com/jinzhu/gorm"
)

type EventRepoPostgreImpl struct {
	db *gorm.DB
}

func CreateEventRepoPostgreImpl(db *gorm.DB) event.EventRepo {
	return &EventRepoPostgreImpl{db}
}

func (e *EventRepoPostgreImpl) BeginTrans() *gorm.DB {
	return e.db.Begin()
}

func (e *EventRepoPostgreImpl) AddEvent(event *models.Event, tx *gorm.DB) error {
	if err := tx.Table("events").Save(&event).Error; err != nil {
		return err
	}
	return nil
}

func (e *EventRepoPostgreImpl) GetAllEvent() (*[]models.Event, error) {
	events := []models.Event{}
	if err := e.db.Preload("Banner").Find(&events).Error; err != nil {
		return nil,err
	}
	return &events,nil
}

