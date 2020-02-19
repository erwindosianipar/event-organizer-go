package banner

import "eventorganizer/golang/models"

type BannerService interface {
	GetAllBanner ()(*[]models.Banner,error)
}
