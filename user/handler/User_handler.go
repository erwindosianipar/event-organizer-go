package handler

import (
	"encoding/json"
	"eventorganizer/golang/models"
	"eventorganizer/golang/user"
	"eventorganizer/golang/util"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type UserHandler struct {
	UserService user.UserService
}

func CreateUserHandler(r *mux.Router, userService user.UserService) {
	userHandler := UserHandler{userService}

	r.HandleFunc("/user/register", userHandler.userRegister).Methods(http.MethodPost)
	r.HandleFunc("/user/login", userHandler.userLogin).Methods(http.MethodPost)
	r.HandleFunc("/user", userHandler.getAllUser).Methods(http.MethodGet)
	r.HandleFunc("/user/{id}", userHandler.getUserByID).Methods(http.MethodGet)
	r.HandleFunc("/user/{id}", userHandler.deleteUser).Methods(http.MethodDelete)
	r.HandleFunc("/user/upgrade", userHandler.upgradeUser).Methods(http.MethodPost)
}

func (h *UserHandler) userRegister(res http.ResponseWriter, req *http.Request) {
	dataRegister, err := ioutil.ReadAll(req.Body)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "Oops, something went wrong.")
		logrus.Error(err)
		return
	}

	reqUser := models.User{}

	err = json.Unmarshal(dataRegister, &reqUser)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "Request body is not valid.")
		logrus.Error(err)
		return
	}

	if !(util.IsEmailValid(reqUser.Email)) {
		util.HandleError(res, http.StatusBadRequest, "Email cannot be empty and must be a valid email.")
		return
	}

	if len(reqUser.Password) < 5 {
		util.HandleError(res, http.StatusBadRequest, "Password cannot be empty and must be more than 5 characters.")
		return
	}

	if len(reqUser.Name) < 3 {
		util.HandleError(res, http.StatusBadRequest, "Name cannot be empty and must be more than 3 characters.")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(reqUser.Password), 10)
	if err != nil {
		logrus.Fatal(err)
	}

	reqUser.Password = string(hash)

	newUser, err := h.UserService.Register(&reqUser)
	if err != nil {
		logrus.Fatal(err)
	}

	if newUser != nil {
		util.HandleSuccess(res, http.StatusCreated, newUser)
	} else {
		util.HandleError(res, http.StatusBadRequest, "Email already used by someone.")
	}
}

func (h *UserHandler) userLogin(res http.ResponseWriter, req *http.Request) {
	dataLogin, err := ioutil.ReadAll(req.Body)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "Oops, something went wrong.")
		logrus.Error(err)
		return
	}

	reqUser := models.User{}

	err = json.Unmarshal(dataLogin, &reqUser)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "Request body is not valid.")
		logrus.Error(err)
		return
	}

	if !(util.IsEmailValid(reqUser.Email)) {
		util.HandleError(res, http.StatusBadRequest, "Email cannot be empty and must be a valid email.")
		return
	}

	if len(reqUser.Password) < 5 {
		util.HandleError(res, http.StatusBadRequest, "Password cannot be empty and must be more than 5 characters.")
		return
	}

	if !(h.UserService.IsAnyEmailUser(reqUser.Email)) {
		util.HandleError(res, http.StatusForbidden, "Email is not match with any email in our database.")
		return
	}

	dataUser, err := h.UserService.GetUserByEmail(reqUser.Email)
	if err != nil {
		logrus.Error(err)
		return
	}

	inputPassword := []byte(reqUser.Password)
	hashedPassword := dataUser.Password

	dataUser = &models.User{
		OrmModel: models.OrmModel{
			ID: dataUser.OrmModel.ID,
		},
		Email:            dataUser.Email,
		Name:             dataUser.Name,
		Role:             dataUser.Role,
		SubmissionStatus: dataUser.SubmissionStatus,
		EventOrganizer: models.EventOrganizer{
			NameEo:     dataUser.EventOrganizer.NameEo,
			KTPNumber:  dataUser.EventOrganizer.KTPNumber,
			KTPPhoto:   dataUser.EventOrganizer.KTPPhoto,
			SIUPNumber: dataUser.EventOrganizer.SIUPNumber,
			IsVerify:   dataUser.EventOrganizer.IsVerify,
		},
	}

	if util.IsPasswordSame(hashedPassword, inputPassword) {
		util.HandleSuccess(res, http.StatusOK, dataUser)
		return
	} else {
		util.HandleError(res, http.StatusForbidden, "Email or password is not match.")
		return
	}
}

func (h *UserHandler) getAllUser(res http.ResponseWriter, req *http.Request) {
	listUser, err := h.UserService.GetAllUser()
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "Oops, something went wrong.")
		return
	}

	util.HandleSuccess(res, http.StatusOK, listUser)
}

func (h *UserHandler) getUserByID(res http.ResponseWriter, req *http.Request) {
	param := mux.Vars(req)
	id, err := strconv.Atoi(param["id"])
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "Please provide valid user id.")
		return
	}

	response, err := h.UserService.GetUserByID(id)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, err.Error())
		return
	}

	util.HandleSuccess(res, http.StatusOK, response)
}

func (h *UserHandler) deleteUser(res http.ResponseWriter, req *http.Request) {
	param := mux.Vars(req)
	id, err := strconv.Atoi(param["id"])
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "Please provide valid user id.")
		return
	}

	_, err = h.UserService.GetUserByID(id)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "No data user with id you entered.")
		return
	}

	_, err = h.UserService.DeleteUser(id)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, err.Error())
		return
	}

	util.HandleError(res, http.StatusOK, "User has been deleted.")
}

func (h *UserHandler) upgradeUser(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if len(id) < 1 {
		util.HandleError(res, http.StatusBadRequest, "id cannot be empty.")
		return
	}

	ID, err := strconv.Atoi(id)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "id must be a valid number.")
		return
	}

	_, err = h.UserService.GetUserByID(ID)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "No data user with ID you entered.")
		return
	}

	nameEO := req.FormValue("name_eo")
	if len(nameEO) < 5 {
		util.HandleError(res, http.StatusBadRequest, "name_eo cannot be empty and must more than 5 characters.")
		return
	}

	ktpNumber := req.FormValue("ktp_number")
	if len(ktpNumber) < 10 {
		util.HandleError(res, http.StatusBadRequest, "ktp_number cannot be empty and must more than 10 characters.")
		return
	}

	siupNumber := req.FormValue("siup_number")
	if len(siupNumber) < 10 {
		util.HandleError(res, http.StatusBadRequest, "siup_number cannot be empty and must more than 10 characters.")
		return
	}

	ktpPhoto, handler, err := req.FormFile("ktp_photo")
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "ktp_photo must be a picture file and cannot be empty.")
		logrus.Error(err)
		return
	}
	defer ktpPhoto.Close()

	random := randstr.Hex(10)
	fileName := fmt.Sprintf("%s%s", "ktp-user-"+id+"-"+random, filepath.Ext(handler.Filename))

	dir, err := os.Getwd()
	if err != nil {
		util.HandleError(res, http.StatusInternalServerError, "Oops, something went wrong.")
		logrus.Error(err)
		return
	}

	fileLocation := filepath.Join(dir, "assets", fileName)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		util.HandleError(res, http.StatusInternalServerError, "Oops, something went wrong.")
		logrus.Error(err)
		return
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, ktpPhoto); err != nil {
		util.HandleError(res, http.StatusInternalServerError, "Oops, something went wrong.")
		logrus.Error(err)
		return
	}

	reqUser := models.User{
		OrmModel:         models.OrmModel{
			ID: uint(ID),
		},
		EventOrganizer: models.EventOrganizer{
			NameEo:     nameEO,
			KTPNumber:  ktpNumber,
			KTPPhoto:   fileName,
			SIUPNumber: siupNumber,
		},
	}

	newUser, err := h.UserService.UpgradeUser(&reqUser)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "Oops, something went wrong.")
		return
	}

	util.HandleSuccess(res, http.StatusCreated, newUser)
}
