package controllers

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/config"
	"go-rest-api/models"
	"net/http"
)

// GetUsers godoc
// @Summary      List all users
// @Description  Get all users
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      404  {object}  nil
// @Router       /users [get]
func GetUsers(c *gin.Context) {
	var users []models.User
	err := config.DB.Find(&users).Error

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}

	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary      Get a user
// @Description  Get a user by id
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      404  {object}  nil
// @Router       /users/{id} [get]
func GetUser(c *gin.Context) {
	var user models.User

	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary      Create a user
// @Description  Create a new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      201  {object}  models.User
// @Failure      400  {object}  nil
// @Failure      500  {object}  nil
// @Router       /users [post]
func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var organization models.Organization
	if err := config.DB.Where("id = ?", user.OrganizationID).First(&organization).Error; err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete a user by id
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  nil
// @Failure      404  {object}  nil
// @Router       /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	config.DB.Delete(&user)

	c.Status(http.StatusNoContent)
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Update user info
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      404  {object}  nil
// @Failure      400  {object}  nil
// @Router       /users/{id} [patch]
func UpdateUser(c *gin.Context) {
	var user models.User
	if err := config.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
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
		if err := config.DB.Where("id = ?", input.OrganizationID).First(&organization).Error; err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}

	input.ID = user.ID
	config.DB.Model(&user).Updates(input)

	c.JSON(http.StatusOK, user)
}

// GetUsersByOrganization godoc
// @Summary      List all users in the given organization
// @Description  Get all users in the given organization
// @Tags         user
// @Tags         organization
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      404  {object}  nil
// @Router       /organizations/{id}/users [get]
func GetUsersByOrganization(c *gin.Context) {
	var organization models.Organization
	var users []models.User

	if err := config.DB.Where("id = ?", c.Param("id")).First(&organization).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err := config.DB.Model(&organization).Association("Users").Find(&users); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
