package guestcontroller

import (
	guestService "regroup-api/services"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	GuestService *guestService.GuestService
}

func (c *Controller) SaveGuest(gin *gin.Context) {
	c.GuestService.SaveGuest(gin)
}

func (c *Controller) GetGuests(gin *gin.Context) {
	c.GuestService.GetGuests(gin)
}

func (c *Controller) GetGuest(gin *gin.Context) {
	c.GuestService.GetGuest(gin)
}

func (c *Controller) DeleteGuest(gin *gin.Context) {
	c.GuestService.DeleteGuest(gin)
}

func (c *Controller) UpdateGuest(gin *gin.Context) {
	c.GuestService.UpdateGuest(gin)
}

func (c *Controller) RegisterRoutes(router *gin.Engine) {
	router.GET("/guests", c.GetGuests)
	router.GET("/guests/:id", c.GetGuest)
	router.PUT("/guests/:id", c.UpdateGuest)
	router.POST("/guests", c.SaveGuest)
	router.DELETE("/guests/:id", c.DeleteGuest)
}
