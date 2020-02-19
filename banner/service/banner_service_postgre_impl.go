package service

import (
	"eventorganizer/golang/banner"
	"eventorganizer/golang/models"
)

type BannerServicePostgreImpl struct {
	bannerRepo banner.BannerRepo
}

func CreateBannerServicePostgreImpl(bannerRepo banner.BannerRepo) banner.BannerService {
	return &BannerServicePostgreImpl{bannerRepo}
}

func (b *BannerServicePostgreImpl) GetAllBanner() (*[]models.Banner, error) {
	return b.bannerRepo.GetAllBanner()
}



