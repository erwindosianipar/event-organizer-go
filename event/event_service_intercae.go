package event

import (
	"eventorganizer/golang/models"
)

type EventService interface {
	TransactionsEvent (event *models.Event,fileName []string)error
	GetAllEvent()(*[]models.Event,error)
}
