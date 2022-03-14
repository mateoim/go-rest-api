package controllers

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/config"
	"go-rest-api/models"
	"net/http"
)

// GetOrganizations godoc
// @Summary      List all organizations
// @Description  Get all organizations
// @Tags         organization
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Organization
// @Failure      404  {object}  nil
// @Router       /organizations [get]
func GetOrganizations(c *gin.Context) {
	var organizations []models.Organization
	err := config.DB.Find(&organizations).Error

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, organizations)
	}
}

// CreateOrganization godoc
// @Summary      Create an organization
// @Description  Create a new organization
// @Tags         organization
// @Accept       json
// @Produce      json
// @Success      201  {object}  models.Organization
// @Failure      400  {object}  nil
// @Failure      500  {object}  nil
// @Router       /organizations [post]
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

// GetOrganization godoc
// @Summary      Get an organization
// @Description  Get an organization by id
// @Tags         organization
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Organization
// @Failure      404  {object}  nil
// @Router       /organizations/{id} [get]
func GetOrganization(c *gin.Context) {
	var organization models.Organization

	if err := config.DB.Where("id = ?", c.Param("id")).First(&organization).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, organization)
}

// DeleteOrganization godoc
// @Summary      Delete an organization
// @Description  Delete an organization by id
// @Tags         organization
// @Accept       json
// @Produce      json
// @Success      200  {object}  nil
// @Failure      404  {object}  nil
// @Router       /organizations/{id} [delete]
func DeleteOrganization(c *gin.Context) {
	var organization models.Organization
	if err := config.DB.Where("id = ?", c.Param("id")).First(&organization).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	config.DB.Delete(&organization)

	c.Status(http.StatusNoContent)
}

// UpdateOrganization godoc
// @Summary      Update an organization
// @Description  Update organization info
// @Tags         organization
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Organization
// @Failure      404  {object}  nil
// @Failure      400  {object}  nil
// @Router       /organizations/{id} [patch]
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
