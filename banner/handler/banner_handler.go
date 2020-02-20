package handler

import (
	"eventorganizer/golang/banner"
	"eventorganizer/golang/util"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type BannerHandler struct {
	bannerService banner.BannerService
}

func CreateBannerHandler(r *mux.Router, bannerService banner.BannerService) {
	bannerHandler := BannerHandler{bannerService}

	r.HandleFunc("/banner", bannerHandler.GetAllBanner).Methods(http.MethodGet)
}

func (b *BannerHandler) GetAllBanner(res http.ResponseWriter, req *http.Request) {
	listBanner,err := b.bannerService.GetAllBanner()
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, err.Error())
		logrus.Error(err)
		return
	}

	util.HandleSuccess(res, http.StatusOK, listBanner)
}
