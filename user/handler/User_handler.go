package handler

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strconv"

	"eventorganizer/golang/models"
	"eventorganizer/golang/user"
	"eventorganizer/golang/util"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
		util.HandleError(res, http.StatusBadRequest,"Name cannot be empty and must be more than 3 characters.")
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
		util.HandleError(res, http.StatusBadRequest,"Email already used by someone.")
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
		util.HandleError(res, http.StatusBadRequest,"Email cannot be empty and must be a valid email.")
		return
	}

	if len(reqUser.Password) < 5 {
		util.HandleError(res, http.StatusBadRequest,"Password cannot be empty and must be more than 5 characters.")
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
		OrmModel:       models.OrmModel{
			ID: dataUser.OrmModel.ID,
		},
		Email:          dataUser.Email,
		Name:           dataUser.Name,
		Role:           dataUser.Role,
		EventOrganizer: models.EventOrganizer{
			IsVerify: false,
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
	dataUser, err := ioutil.ReadAll(req.Body)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "Oops, something went wrong.")
		logrus.Error(err)
		return
	}

	reqUser := models.User{}

	err = json.Unmarshal(dataUser, &reqUser)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "Request body is not valid.")
		logrus.Error(err)
		return
	}

	if reqUser.ID == 0 {
		util.HandleError(res, http.StatusBadRequest, "ID user cannot be empty.")
		return
	}

	if len(reqUser.EventOrganizer.NameEo) < 5 {
		util.HandleError(res, http.StatusBadRequest, "EO name cannot be empty and must more than 5 characters.")
		return
	}

	if len(reqUser.EventOrganizer.KTPNumber) < 10 {
		util.HandleError(res, http.StatusBadRequest, "KTP number cannot be empty and must more than 10 characters.")
		return
	}

	if reqUser.EventOrganizer.KTPPhoto == "" {
		util.HandleError(res, http.StatusBadRequest, "KTP photo cannot be empty.")
		return
	}

	if len(reqUser.EventOrganizer.SIUPNumber) < 10 {
		util.HandleError(res, http.StatusBadRequest, "SIUP number cannot be empty and must more than 10 characters.")
		return
	}

	id := int(reqUser.ID)

	_, err = h.UserService.GetUserByID(id)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "No data user with id you entered.")
		return
	}

	newUser, err := h.UserService.UpgradeUser(&reqUser)
	if err != nil {
		util.HandleError(res, http.StatusBadRequest, "Oops, something went wrong.")
		return
	}

	util.HandleSuccess(res, http.StatusCreated, newUser)
}