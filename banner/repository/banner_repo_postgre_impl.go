package repository

import (
	"eventorganizer/golang/banner"
	"eventorganizer/golang/models"
	"github.com/jinzhu/gorm"
)

type BannerRepoPostgreImpl struct {
	db *gorm.DB
}

func CreateBannerRepoPostgreImpl(db *gorm.DB) banner.BannerRepo {
	return &BannerRepoPostgreImpl{db}
}

func (b *BannerRepoPostgreImpl) AddBanner(banner *models.Banner, tx *gorm.DB) error {
	if err := tx.Table("banners").Save(&banner).Error; err != nil {
		return err
	}
	return nil
}

func (b *BannerRepoPostgreImpl) GetAllBanner() (*[]models.Banner, error) {
	banners := []models.Banner{}
	if err := b.db.Table("banners").Find(&banners).Error; err != nil {
		return nil,err
	}
	return &banners,nil
}

func (b *BannerRepoPostgreImpl) GettBannerByIdEvent(id_event int) (*[]models.Banner, error) {
	banners:=[]models.Banner{}
	if err:= b.db.Where("id_event = ?",id_event).Find(&banners).Error;err!=nil{
		return nil,err
	}
	return &banners,nil
}