package service

import (
	"eventorganizer/golang/banner"
	"eventorganizer/golang/event"
	"eventorganizer/golang/models"
)

type EventServicePostgreImpl struct {
	eventRepo event.EventRepo
	bannerRepo banner.BannerRepo
}

func CreateEventServicePostgreImpl(eventRepo event.EventRepo, bannerRepo banner.BannerRepo) event.EventService {
	return &EventServicePostgreImpl{eventRepo,bannerRepo}
}

func (e *EventServicePostgreImpl) TransactionsEvent(event *models.Event, fileName []string) error {
	tx := e.eventRepo.BeginTrans()
	err := e.eventRepo.AddEvent(event,tx)
	if err!=nil{
		tx.Rollback()
		return err
	}

	//id_event := event.ID
	//for _,name := range fileName{
	//	var yourBanner  = models.Banner{ID_Event:id_event,Banner_Foto:name}
	//	err = e.bannerRepo.AddBanner(&yourBanner,tx)
	//	fmt.Println(name)
	//	if err!=nil{
	//		tx.Rollback()
	//		return err
	//	}
	//}
	return tx.Commit().Error
}

func (e *EventServicePostgreImpl) GetAllEvent() (*[]models.Event, error) {
	return e.eventRepo.GetAllEvent()
}
