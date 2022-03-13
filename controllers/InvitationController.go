package controllers

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/config"
	"go-rest-api/models"
	"net/http"
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

	c.JSON(http.StatusNoContent, nil)
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
		invitation.Status = models.Accepted
	} else {
		invitation.Status = models.Rejected
	}

	config.DB.Model(&invitation).Updates(invitation)

	c.JSON(http.StatusOK, invitation)
}
