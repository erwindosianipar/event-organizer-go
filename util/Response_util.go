package util

import (
	"encoding/json"
	"net/http"

	"eventorganizer/golang/models"

	"github.com/sirupsen/logrus"
)

func HandleSuccess(res http.ResponseWriter, status int, data interface{})  {
	returnData := models.ResponseWrapper{
		Success: true,
		Message: "success",
		Data:    data,
	}

	jsonData, err := json.Marshal(returnData)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		_, _ = res.Write([]byte("Oops, something went wrong."))
		logrus.Error(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	_, _ = res.Write(jsonData)
}

func HandleError(res http.ResponseWriter, status int, message string)  {
	returnData := models.ResponseWrapper{
		Success: false,
		Message: message,
	}

	jsonData, err := json.Marshal(returnData)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		_, _ = res.Write([]byte("Oops, something went wrong."))
		logrus.Error(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	_, _ = res.Write(jsonData)
}
