package controllers

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/models"
	"net/http"
	"strconv"
	"time"
)

func GetEvents(c *gin.Context) {
	var events []models.Event
	err := models.DB.Find(&events).Error

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, events)
	}
}

func CreateEvent(c *gin.Context) {
	var event models.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if time.Time(event.EndDate).Before(time.Time(event.StartDate)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date cannot be before Start date"})
		return
	}

	if err := models.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

func GetEvent(c *gin.Context) {
	var event models.Event

	if err := models.DB.Where("id = ?", c.Param("id")).First(&event).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, event)
}

func DeleteEvent(c *gin.Context) {
	var event models.Event
	if err := models.DB.Where("id = ?", c.Param("id")).First(&event).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	models.DB.Delete(&event)

	c.JSON(http.StatusNoContent, nil)
}

func UpdateEvent(c *gin.Context) {
	var event models.Event
	if err := models.DB.Where("id = ?", c.Param("id")).First(&event).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var input models.Event
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = event.ID
	input.StartDate = event.StartDate
	input.EndDate = event.EndDate
	models.DB.Model(&event).Updates(input)

	c.JSON(http.StatusOK, event)
}

func RegisterUser(c *gin.Context) {
	var event models.Event
	if err := models.DB.Where("id = ?", c.Param("id")).First(&event).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var form models.IDForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := models.DB.Where("id = ?", strconv.Itoa(int(form.ID))).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := models.DB.Model(&event).Association("Users").Append(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
