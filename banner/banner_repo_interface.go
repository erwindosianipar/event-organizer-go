package banner

import (
	"eventorganizer/golang/models"
	"github.com/jinzhu/gorm"
)

type BannerRepo interface {
	AddBanner (banner *models.Banner,tx *gorm.DB)error
	GetAllBanner ()(*[]models.Banner,error)
}
