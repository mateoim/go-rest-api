package controllers

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/models"
	"net/http"
	"time"
)

func GetMeetings(c *gin.Context) {
	var event models.Event
	var meetings []models.Meeting

	if err := models.DB.Where("id = ?", c.Param("id")).First(&event).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := models.DB.Model(&event).Association("Meetings").Find(&meetings); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, meetings)
}

func CreateMeeting(c *gin.Context) {
	var event models.Event
	var meeting models.Meeting

	if err := models.DB.Where("id = ?", c.Param("id")).First(&event).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := c.ShouldBindJSON(&meeting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eventStartTime := time.Time(event.StartDate).UTC()
	eventEndTime := time.Time(event.EndDate).UTC()
	meetingStartTime := time.Time(meeting.StartDate).UTC()
	meetingEndTime := time.Time(meeting.EndDate).UTC()

	if meetingEndTime.Before(meetingStartTime) || meetingStartTime.Before(eventStartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End date cannot be before Start date"})
		return
	} else if meetingEndTime.After(eventEndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Meeting End date cannot be after event End date"})
		return
	}

	if err := models.DB.Create(&meeting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, meeting)
}

func GetMeeting(c *gin.Context) {
	var meeting models.Meeting

	if err := GetMeetingModel(c, &meeting); err != nil {
		return
	}

	c.JSON(http.StatusOK, meeting)
}

func GetMeetingModel(c *gin.Context, meeting *models.Meeting) error {
	var event models.Event

	if err := models.DB.Where("id = ?", c.Param("id")).First(&event).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return err
	}

	if err := models.DB.Where("id = ?", c.Param("meeting-id")).First(&meeting).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return err
	}

	return nil
}

func DeleteMeeting(c *gin.Context) {
	var meeting models.Meeting

	if err := GetMeetingModel(c, &meeting); err != nil {
		return
	}

	models.DB.Delete(&meeting)

	c.JSON(http.StatusNoContent, nil)
}