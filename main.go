package main

import (
	"database/sql"
	"expvar"
	"log"
	"net/http"
	"runtime"
	"time"

	"gins/config"
	"gins/controller"
	"gins/helper"
	"gins/model"
	"gins/repository"
	"gins/service"

	"gins/router"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func main() {

	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	// Database
	db := config.ConnectionDB(&loadConfig)
	validate := validator.New()

	db.Table("tags").AutoMigrate(&model.Tags{})
	db.Table("users").AutoMigrate(&model.Users{})

	// Initialize expvar metrics
	expvar.NewString("version").Set("1.0.0")
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() interface{} {
		return getDBStats(db)
	}))
	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	// Repository
	tagRespository := repository.NewTagsRepositoryImpl(db)
	userRepository := repository.NewUsersRepositoryImpl(db)

	// Service
	tagsService := service.NewTagsServiceImpl(&tagRespository, validate)
	authenticationService := service.NewAuthenticationServiceImpl(userRepository, validate)

	// Controller
	tagController := controller.NewTagsController(tagsService)
	authenticationController := controller.NewAuthenticationController(authenticationService)
	usersController := controller.NewUsersController(userRepository)

	// Router
	router := router.NewRouter(userRepository, tagController, authenticationController, usersController)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err = server.ListenAndServe()
	helper.ErrorPanic(err)
}

func getDBStats(db *gorm.DB) sql.DBStats {
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("error getting DB object:", err)
		return sql.DBStats{}
	}
	return sqlDB.Stats()
}
