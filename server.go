package main

import (
	guestcontroller "regroup-api/controllers"
	"regroup-api/database"
	guestrepository "regroup-api/repositories"
	guestservice "regroup-api/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitGuestController(db *mongo.Client, router *gin.Engine) {
	guestRepository := guestrepository.GuestRepository{DB: db, Collection: database.GetCollection(db, "guests")}
	guestService := guestservice.GuestService{Repository: &guestRepository}
	guestController := guestcontroller.Controller{GuestService: &guestService}
	guestController.RegisterRoutes(router)
}

func InitHealthCheck(router *gin.Engine) {
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})
}

func InitControllers(db *mongo.Client, router *gin.Engine) {
	InitGuestController(db, router)
	InitHealthCheck(router)
}

func main() {
	router := gin.Default()
	db := database.Init()

	InitControllers(db, router)

	router.Run("0.0.0.0:8080")
}
