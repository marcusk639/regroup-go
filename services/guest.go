package guestservice

import (
	"log"
	"net/http"
	models "regroup-api/models"
	guestrepository "regroup-api/repositories"

	"github.com/gin-gonic/gin"
)

type GuestService struct {
	Repository *guestrepository.GuestRepository
}

func (service *GuestService) SaveGuest(c *gin.Context) {
	guest := models.Guest{}

	if err := c.ShouldBindJSON(&guest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result, err := service.Repository.SaveGuest(guest)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error saving guest"})
		return
	}

	c.JSON(http.StatusOK, *result)
}

func (service *GuestService) GetGuests(c *gin.Context) {
	guests, err := service.Repository.GetGuests()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, *guests)
}

func (service *GuestService) GetGuest(c *gin.Context) {
	id := c.Param("id")

	guest, err := service.Repository.GetGuest(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, *guest)
}

func (service *GuestService) DeleteGuest(c *gin.Context) {
	id := c.Param("id")

	err := service.Repository.DeleteGuest(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error deleting guest"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Guest deleted"})
}

func (service *GuestService) UpdateGuest(c *gin.Context) {
	id := c.Param("id")

	var guest models.Guest

	if err := c.ShouldBindJSON(&guest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.Repository.UpdateGuest(id, guest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error updating guest"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Guest updated"})
}
