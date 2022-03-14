package controllers

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/config"
	"go-rest-api/models"
	"net/http"
	"strconv"
	"time"
)

// GetEvents godoc
// @Summary      List all events
// @Description  Get all events
// @Tags         event
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Event
// @Failure      404  {object}  nil
// @Router       /events [get]
func GetEvents(c *gin.Context) {
	var events []models.Event
	err := config.DB.Find(&events).Error

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, events)
	}
}

// CreateEvent godoc
// @Summary      Create an event
// @Description  Create a new event
// @Tags         event
// @Accept       json
// @Produce      json
// @Success      201  {object}  models.Event
// @Failure      400  {object}  nil
// @Failure      500  {object}  nil
// @Router       /events [post]
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

	if err := config.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetEvent godoc
// @Summary      Get an event
// @Description  Get an event by id
// @Tags         event
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Event
// @Failure      404  {object}  nil
// @Router       /events/{id} [get]
func GetEvent(c *gin.Context) {
	var event models.Event

	if err := config.DB.Where("id = ?", c.Param("id")).First(&event).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, event)
}

// DeleteEvent godoc
// @Summary      Delete an event
// @Description  Delete an event by id
// @Tags         event
// @Accept       json
// @Produce      json
// @Success      200  {object}  nil
// @Failure      404  {object}  nil
// @Router       /event/{id} [delete]
func DeleteEvent(c *gin.Context) {
	var event models.Event
	if err := config.DB.Where("id = ?", c.Param("id")).First(&event).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	config.DB.Delete(&event)

	c.Status(http.StatusNoContent)
}

// UpdateEvent godoc
// @Summary      Update an event
// @Description  Update event info
// @Tags         event
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Event
// @Failure      404  {object}  nil
// @Failure      400  {object}  nil
// @Router       /events/{id} [patch]
func UpdateEvent(c *gin.Context) {
	var event models.Event
	if err := config.DB.Where("id = ?", c.Param("id")).First(&event).Error; err != nil {
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
	config.DB.Model(&event).Updates(input)

	c.JSON(http.StatusOK, event)
}

// RegisterUser godoc
// @Summary      Registers a user to the given event
// @Description  Registers a user to the given event by sending user id in JSON body
// @Tags         event
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      201  {object}  nil
// @Failure      404  {object}  nil
// @Failure      400  {object}  nil
// @Failure      500  {object}  nil
// @Router       /events/{id}/register [post]
func RegisterUser(c *gin.Context) {
	var event models.Event
	if err := config.DB.Where("id = ?", c.Param("id")).First(&event).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var form models.IDForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("id = ?", strconv.Itoa(int(form.ID))).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := config.DB.Model(&event).Association("Users").Append(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
