package handler

import (
	"eventorganizer/golang/event"
	"eventorganizer/golang/models"
	"eventorganizer/golang/util"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/thanhpk/randstr"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type EventHandler struct {
	eventService event.EventService
}

func CreateEventHandler(r *mux.Router, eventService event.EventService) {
	userHandler := EventHandler{eventService}

	r.HandleFunc("/event", userHandler.TransactionsEvent).Methods(http.MethodPost)
	r.HandleFunc("/event", userHandler.GetAllEvent).Methods(http.MethodGet)
}
func (e *EventHandler) TransactionsEvent(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id_user")
	if len(id) < 1 {
		util.HandleError(res, http.StatusBadRequest, "id cannot be empty.")
		return
	}

	ID, err := strconv.Atoi(id)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "id must be a valid number.")
		return
	}

	reqBanner := []string{}
	req.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile
	fhs := req.MultipartForm.File["banner_foto"]
	for _, fh := range fhs {
		f, err := fh.Open()
		random := randstr.Hex(5)
		fileName := fmt.Sprintf("%s%s", "banner-foto-"+id+"-"+random, filepath.Ext(fh.Filename))
		dir, err := os.Getwd()
		if err != nil {
			util.HandleError(res, http.StatusInternalServerError, "Oops, something went wrong.")
			logrus.Error(err)
			return
		}
		fileLocation := filepath.Join(dir, "assets", fileName)
		reqBanner = append(reqBanner,fileName)
		targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			util.HandleError(res, http.StatusInternalServerError, "Oops, something went wrong.")
			logrus.Error(err)
			return
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, f); err != nil {
			util.HandleError(res, http.StatusInternalServerError, "Oops, something went wrong.")
			logrus.Error(err)
			return
		}
		f.Close()
	}

	name := req.FormValue("name")
	if name == "" {
		util.HandleError(res, http.StatusBadRequest, "name event cannot be empty.")
		return
	}
	lokasi := req.FormValue("lokasi")
	if lokasi == "" {
		util.HandleError(res, http.StatusBadRequest, "location event cannot be empty.")
		return
	}
	event_date := req.FormValue("event_date")
	if event_date == "" {
		util.HandleError(res, http.StatusBadRequest, "event_date event cannot be empty.")
		return
	}
	kuota := req.FormValue("kuota")
	if kuota == "" {
		util.HandleError(res, http.StatusBadRequest, "quota event cannot be empty.")
		return
	}
	yourKuota, err := strconv.Atoi(kuota)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "quota must be a valid number.")
		return
	}
	harga := req.FormValue("harga")
	yourHarga, err := strconv.Atoi(harga)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "price must be a valid number.")
		return
	}

	var yourBanner  = []models.Banner{}
	for _,name := range reqBanner{
		yourBanner = append(yourBanner,models.Banner{Banner_Foto:name})
	}
	reqEvent := models.Event{
		ID_User:    uint(ID),
		Name:       name,
		Lokasi:     lokasi,
		Event_date: event_date,
		Kuota:      yourKuota,
		Harga:      yourHarga,
		Banner:yourBanner,
	}


	err = e.eventService.TransactionsEvent(&reqEvent, reqBanner)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, err.Error())
		logrus.Error(err)
		return
	}

	util.HandleSuccess(res, http.StatusCreated, reqEvent)
	return
}

func (e *EventHandler) GetAllEvent(res http.ResponseWriter, req *http.Request){
	listEvent,err:=e.eventService.GetAllEvent()
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, err.Error())
		logrus.Error(err)
		return
	}

	util.HandleSuccess(res, http.StatusOK, listEvent)
}