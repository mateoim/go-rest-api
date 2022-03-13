package controllers

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/config"
	"go-rest-api/models"
	"net/http"
	"strconv"
	"time"
)

func GetInvitations(c *gin.Context) {
	var meeting models.Meeting
	var invitations []models.Invitation

	if err := GetMeetingModel(c, &meeting); err != nil {
		return
	}

	if err := config.DB.Model(&meeting).Association("Invitations").Find(&invitations); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, invitations)
}

func CreateInvitation(c *gin.Context) {
	var meeting models.Meeting
	var invitation models.Invitation

	if err := GetMeetingModel(c, &meeting); err != nil {
		return
	}

	if err := c.ShouldBindJSON(&invitation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("meeting-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invitation.MeetingID = uint(id)

	var user models.User
	if err := config.DB.Where("id = ?", invitation.UserID).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var userEvents []models.Event
	if err := config.DB.Model(&user).Association("Events").Find(&userEvents); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	eventFound := false

	for _, event := range userEvents {
		if event.ID == meeting.EventID {
			eventFound = true
			break
		}
	}

	if !eventFound {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User is not registered to this event."})
		return
	}

	invitation.Status = models.Pending
	if err := config.DB.Create(&invitation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, invitation)
}

func GetInvitation(c *gin.Context) {
	var invitation models.Invitation

	if err := GetInvitationModel(c, &invitation); err != nil {
		return
	}

	c.JSON(http.StatusOK, invitation)
}

func GetInvitationModel(c *gin.Context, invitation *models.Invitation) error {
	var meeting models.Meeting

	if err := GetMeetingModel(c, &meeting); err != nil {
		return err
	}

	if err := config.DB.Where("id = ?", c.Param("invitation-id")).First(&invitation).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return err
	}

	return nil
}

func DeleteInvitation(c *gin.Context) {
	var invitation models.Invitation

	if err := GetInvitationModel(c, &invitation); err != nil {
		return
	}

	config.DB.Delete(&invitation)

	c.Status(http.StatusNoContent)
}

func Accept(c *gin.Context) {
	UpdateStatus(c, true)
}

func Reject(c *gin.Context) {
	UpdateStatus(c, false)
}

func UpdateStatus(c *gin.Context, accepted bool) {
	var invitation models.Invitation

	if err := GetInvitationModel(c, &invitation); err != nil {
		return
	}

	if accepted {
		if CheckConflicts(c, &invitation) {
			c.JSON(http.StatusConflict, gin.H{"error": "Cannot accept a meeting invitation in conflict with a meeting."})
			return
		}
		invitation.Status = models.Accepted
	} else {
		invitation.Status = models.Rejected
	}

	config.DB.Model(&invitation).Updates(invitation)

	c.JSON(http.StatusOK, invitation)
}

func CheckConflicts(c *gin.Context, invitation *models.Invitation) bool {
	var user models.User
	var meetings []models.Meeting

	if err := config.DB.Where("id = ?", invitation.UserID).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return true
	}

	var meeting models.Meeting
	if err := config.DB.Where("id = ?", invitation.MeetingID).First(&meeting).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return true
	}

	meetingStartTime := time.Time(meeting.StartDate).UTC()
	meetingEndTime := time.Time(meeting.EndDate).UTC()

	if err := config.DB.Model(&user).Association("Meetings").Find(&meetings); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return true
	}

	for _, m := range meetings {
		if m.ID == invitation.MeetingID {
			continue
		}

		currentStartTime := time.Time(m.StartDate).UTC()
		currentEndTime := time.Time(m.EndDate).UTC()

		if meetingStartTime.Before(currentEndTime) && currentStartTime.Before(meetingEndTime) {
			return true
		}
	}

	err := config.DB.Model(&user).Association("Meetings").Append(&meeting)
	if err != nil {
		return true
	}

	return false
}
