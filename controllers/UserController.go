package controllers

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/models"
	"net/http"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	err := models.DB.Find(&users).Error

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	var user models.User

	if err := models.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var organization models.Organization
	if err := models.DB.Where("id = ?", user.OrganizationID).First(&organization).Error; err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func DeleteUser(c *gin.Context) {
	var user models.User
	if err := models.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	models.DB.Delete(&user)

	c.JSON(http.StatusNoContent, nil)
}

func UpdateUser(c *gin.Context) {
	var user models.User
	if err := models.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.OrganizationID != 0 && user.OrganizationID != input.OrganizationID {
		var organization models.Organization
		if err := models.DB.Where("id = ?", input.OrganizationID).First(&organization).Error; err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}

	input.ID = user.ID
	models.DB.Model(&user).Updates(input)

	c.JSON(http.StatusOK, user)
}

func GetUsersByOrganization(c *gin.Context) {
	var organization models.Organization
	var users []models.User

	if err := models.DB.Where("id = ?", c.Param("id")).First(&organization).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := models.DB.Model(&organization).Association("Users").Find(&users); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
