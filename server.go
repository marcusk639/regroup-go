package main

import (
	guestcontroller "regroup-api/controllers"
	"regroup-api/database"
	guestservice "regroup-api/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitGuestController(db *mongo.Client, router *gin.Engine) {
	guestService := guestservice.GuestService{DB: db}
	guestController := guestcontroller.Controller{GuestService: &guestService}
	guestController.RegisterRoutes(router)
}

func InitControllers(db *mongo.Client, router *gin.Engine) {
	InitGuestController(db, router)
}

func main() {
	router := gin.Default()
	db := database.Init()

	InitControllers(db, router)

	router.Run("localhost:8080")
}
