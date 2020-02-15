package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"eventorganizer/golang/models"
	userHandler "eventorganizer/golang/user/handler"
	userRepo "eventorganizer/golang/user/repo"
	userService "eventorganizer/golang/user/service"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func init() {
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUser := viper.GetString("database.user")
	dbPass := viper.GetString("database.pass")
	dbName := viper.GetString("database.name")

	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	val := url.Values{}
	val.Add("sslmode", "disable")
	connStr := fmt.Sprintf("%s?%s", connection, val.Encode())

	dbConn, err := gorm.Open("postgres", connStr)
	if err != nil {
		logrus.Fatal(err)
	}

	err = dbConn.DB().Ping()
	if err != nil {
		logrus.Error(err)
	}

	defer func() {
		err = dbConn.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	dbConn.Debug().AutoMigrate(
		&models.User{},
		)

	router := mux.NewRouter().StrictSlash(true)

	userRepo := userRepo.CreateUserRepoPostgreImpl(dbConn)
	userService := userService.CreateUserService(userRepo)
	userHandler.CreateUserHandler(router, userService)

	fmt.Println("Starting web server at port : 8080")
	err = http.ListenAndServe(": " + "8080", router)
	if err != nil {
		logrus.Fatal(err)
	}

}
