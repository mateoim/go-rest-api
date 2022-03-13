package controllers

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/config"
	"go-rest-api/models"
	"net/http"
)

func GetOrganizations(c *gin.Context) {
	var organizations []models.Organization
	err := config.DB.Find(&organizations).Error

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, organizations)
	}
}

func CreateOrganization(c *gin.Context) {
	var organization models.Organization

	if err := c.ShouldBindJSON(&organization); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&organization).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, organization)
}

func GetOrganization(c *gin.Context) {
	var organization models.Organization

	if err := config.DB.Where("id = ?", c.Param("id")).First(&organization).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, organization)
}

func DeleteOrganization(c *gin.Context) {
	var organization models.Organization
	if err := config.DB.Where("id = ?", c.Param("id")).First(&organization).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	config.DB.Delete(&organization)

	c.Status(http.StatusNoContent)
}

func UpdateOrganization(c *gin.Context) {
	var organization models.Organization
	if err := config.DB.Where("id = ?", c.Param("id")).First(&organization).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	var input models.Organization
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = organization.ID
	config.DB.Model(&organization).Updates(input)

	c.JSON(http.StatusOK, organization)
}
