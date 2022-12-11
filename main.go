package main

import (
	"Bananza/db"
	"Bananza/handler"
	"Bananza/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title bananza
// @version 1.0
// description API server for todo application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

var router *gin.Engine

func main() {
	repos := db.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(Server)
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}
