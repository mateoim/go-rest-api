package controllers

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/config"
	"go-rest-api/models"
	"net/http"
	"strconv"
	"time"
)

// GetMeetings godoc
// @Summary      List all meetings for the given event
// @Description  Get all meetings for the given event
// @Tags         meeting
// @Tags         event
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Meeting
// @Failure      404  {object}  nil
// @Router       /events/{id}/meetings [get]
func GetMeetings(c *gin.Context) {
	var event models.Event
	var meetings []models.Meeting

	if err := config.DB.Where("id = ?", c.Param("id")).First(&event).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := config.DB.Model(&event).Association("Meetings").Find(&meetings); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, meetings)
}

// CreateMeeting godoc
// @Summary      Create a meeting for the given event
// @Description  Create a new meeting for the given event
// @Tags         meeting
// @Tags         event
// @Accept       json
// @Produce      json
// @Success      201  {object}  models.Meeting
// @Failure      404  {object}  nil
// @Failure      400  {object}  nil
// @Failure      500  {object}  nil
// @Router       /events/{id}/meetings [post]
func CreateMeeting(c *gin.Context) {
	var event models.Event
	var meeting models.Meeting

	if err := config.DB.Where("id = ?", c.Param("id")).First(&event).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := c.ShouldBindJSON(&meeting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meeting.EventID = uint(id)

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

	if err := config.DB.Create(&meeting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, meeting)
}

// GetMeeting godoc
// @Summary      Get a meeting for the given event
// @Description  Get a meeting for the given event by meeting id
// @Tags         meeting
// @Tags         event
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Meeting
// @Failure      404  {object}  nil
// @Router       /events/{id}/meetings/{meeting-id} [get]
func GetMeeting(c *gin.Context) {
	var meeting models.Meeting

	if err := GetMeetingModel(c, &meeting); err != nil {
		return
	}

	c.JSON(http.StatusOK, meeting)
}

func GetMeetingModel(c *gin.Context, meeting *models.Meeting) error {
	var event models.Event

	if err := config.DB.Where("id = ?", c.Param("id")).First(&event).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return err
	}

	if err := config.DB.Where("id = ?", c.Param("meeting-id")).First(&meeting).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return err
	}

	return nil
}

// DeleteMeeting godoc
// @Summary      Delete a meeting
// @Description  Delete a meeting for the given event by meeting id
// @Tags         meeting
// @Tags         event
// @Accept       json
// @Produce      json
// @Success      200  {object}  nil
// @Failure      404  {object}  nil
// @Router       /events/{id}/meetings/{meeting-id} [delete]
func DeleteMeeting(c *gin.Context) {
	var meeting models.Meeting

	if err := GetMeetingModel(c, &meeting); err != nil {
		return
	}

	config.DB.Delete(&meeting)

	c.Status(http.StatusNoContent)
}

// Schedule godoc
// @Summary      Schedule the given meeting
// @Description  Schedule the given meeting if no pending or rejected invitations are found
// @Tags         meeting
// @Tags         event
// @Accept       json
// @Produce      json
// @Success      200  {object}  nil
// @Failure      404  {object}  nil
// @Failure      409  {object}  nil
// @Router       /events/{id}/meetings/{meeting-id}/schedule [post]
func Schedule(c *gin.Context) {
	var meeting models.Meeting

	if err := GetMeetingModel(c, &meeting); err != nil {
		return
	}

	var invitations []models.Invitation
	if err := config.DB.Model(&meeting).Association("Invitations").Find(&invitations); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	for _, inv := range invitations {
		if inv.Status != models.Accepted {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Pending or rejected invitation found."})
			return
		}
	}

	c.Status(http.StatusOK)
}
