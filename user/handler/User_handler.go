package handler

import (
	"encoding/json"
	"fmt"
	"github.com/thanhpk/randstr"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"eventorganizer/golang/models"
	"eventorganizer/golang/user"
	"eventorganizer/golang/util"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
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
	r.HandleFunc("/user/upgrade/handle", userHandler.handleUpgrade).Methods(http.MethodPut)
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
		util.HandleError(res, http.StatusInternalServerError, "Oops, something went wrong.")
		logrus.Error(err)
		return
	}

	reqUser.Password = string(hash)

	if h.UserService.IsAnyEmailUser(reqUser.Email) {
		util.HandleError(res, http.StatusBadRequest, "Email already used by someone.")
		logrus.Error(err)
		return
	}

	newUser, err := h.UserService.Register(&reqUser)
	if err != nil {
		util.HandleError(res, http.StatusInternalServerError, err.Error())
		logrus.Fatal(err)
		return
	}

	util.HandleSuccess(res, http.StatusCreated, newUser)
	return
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
		logrus.Error(err)
		return
	}

	if len(reqUser.Password) < 5 {
		util.HandleError(res, http.StatusBadRequest, "Password cannot be empty and must be more than 5 characters.")
		logrus.Error(err)
		return
	}

	if !(h.UserService.IsAnyEmailUser(reqUser.Email)) {
		util.HandleError(res, http.StatusForbidden, "Email is not match with any email in our database.")
		logrus.Error(err)
		return
	}

	dataUser, err := h.UserService.GetUserByEmail(reqUser.Email)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, err.Error())
		logrus.Error(err)
		return
	}

	inputPassword := []byte(reqUser.Password)
	hashedPassword := dataUser.Password

	getDataUser, _ := h.UserService.GetUserByID(int(dataUser.OrmModel.ID))
	if util.IsPasswordSame(hashedPassword, inputPassword) {
		util.HandleSuccess(res, http.StatusOK, getDataUser)
		return
	}

	util.HandleError(res, http.StatusForbidden, "Email or password is not match.")
	return
}

func (h *UserHandler) getAllUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	listUser, err := h.UserService.GetAllUser()
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, err.Error())
		logrus.Error(err)
		return
	}

	util.HandleSuccess(res, http.StatusOK, listUser)
	return
}

func (h *UserHandler) getUserByID(res http.ResponseWriter, req *http.Request) {
	param := mux.Vars(req)
	id, err := strconv.Atoi(param["id"])
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "Please provide valid user id.")
		logrus.Error(err)
		return
	}

	response, err := h.UserService.GetUserByID(id)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, err.Error())
		logrus.Error(err)
		return
	}

	util.HandleSuccess(res, http.StatusOK, response)
	return
}

func (h *UserHandler) deleteUser(res http.ResponseWriter, req *http.Request) {
	param := mux.Vars(req)
	id, err := strconv.Atoi(param["id"])
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "Please provide valid user id.")
		logrus.Error(err)
		return
	}

	_, err = h.UserService.GetUserByID(id)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, err.Error())
		logrus.Error(err)
		return
	}

	_, err = h.UserService.DeleteUser(id)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, err.Error())
		logrus.Error(err)
		return
	}

	util.HandleError(res, http.StatusOK, "User has been deleted.")
	return
}

func (h *UserHandler) upgradeUser(res http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	fmt.Println(id)
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
		util.HandleError(res, http.StatusBadRequest, err.Error())
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

	fmt.Printf("%v %v %v %v \n", id, nameEO, ktpNumber, siupNumber)

	ktpPhoto, handler, err := req.FormFile("ktp_photo")
	fmt.Printf("%v", ktpPhoto)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "ktp_photo must be a picture file and cannot be empty.")
		logrus.Error(err)
		return
	}
	defer ktpPhoto.Close()

	random := randstr.Hex(5)
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
		OrmModel: models.OrmModel{
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
		util.HandleError(res, http.StatusBadRequest, err.Error())
		logrus.Error(err)
		return
	}

	util.HandleSuccess(res, http.StatusCreated, newUser)
	return
}

func (h *UserHandler) handleUpgrade(res http.ResponseWriter, req *http.Request) {
	dataUpgrade, err := ioutil.ReadAll(req.Body)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "Oops, something went wrong.")
		logrus.Error(err)
		return
	}

	reqUser := models.User{}

	err = json.Unmarshal(dataUpgrade, &reqUser)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "Request body is not valid.")
		logrus.Error(err)
		return
	}

	if reqUser.OrmModel.ID == 0 {
		util.HandleError(res, http.StatusBadRequest, "id user cannot be empty.")
		logrus.Error(err)
		return
	}

	_, err = h.UserService.GetUserByID(int(reqUser.OrmModel.ID))
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, err.Error())
		logrus.Error(err)
		return
	}

	status := reqUser.SubmissionStatus
	if !(status == "rejected" || reqUser.SubmissionStatus == "accepted") {
		util.HandleError(res, http.StatusBadRequest, "submission_status must be filled in with rejected or accepted.")
		logrus.Error(err)
		return
	}

	dataUser, err := h.UserService.HandleUpgrade(int(reqUser.OrmModel.ID), status)
	if err != nil {
		util.HandleError(res, http.StatusInternalServerError, err.Error())
		logrus.Error(err)
		return
	}

	util.HandleSuccess(res, http.StatusOK, dataUser)
	return
}

