package main

import (
	bannerHandler "eventorganizer/golang/banner/handler"
	bannerRepo "eventorganizer/golang/banner/repository"
	bannerService "eventorganizer/golang/banner/service"
	eventHandler "eventorganizer/golang/event/handler"
	eventRepo "eventorganizer/golang/event/repository"
	eventService "eventorganizer/golang/event/service"
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
		&models.User{},&models.Event{},&models.Banner{},
		)

	router := mux.NewRouter().StrictSlash(true)

	userRepo := userRepo.CreateUserRepoPostgreImpl(dbConn)
	eventRepo := eventRepo.CreateEventRepoPostgreImpl(dbConn)
	bannerRepo := bannerRepo.CreateBannerRepoPostgreImpl(dbConn)


	userService := userService.CreateUserService(userRepo)
	eventService := eventService.CreateEventServicePostgreImpl(eventRepo,bannerRepo)
	bannerService := bannerService.CreateBannerServicePostgreImpl(bannerRepo)


	userHandler.CreateUserHandler(router, userService)
	eventHandler.CreateEventHandler(router,eventService)
	bannerHandler.CreateBannerHandler(router,bannerService)

	go serverImage()

	fmt.Println("Starting web server at port : 8082")
	err = http.ListenAndServe(": " + "8082", router)
	if err != nil {
		logrus.Fatal(err)
	}

}

func serverImage()  {
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/", fs)

	port := "8083"
	fmt.Printf("Starting image server at http://localhost:%s/\n", port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
