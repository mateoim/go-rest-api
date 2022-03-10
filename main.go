package main

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/controllers"
	"go-rest-api/models"
)

func main() {
	r := gin.Default()
	models.ConnectToDatabase()

	r.GET("/organizations", controllers.GetOrganizations)
	r.POST("/organizations", controllers.CreateOrganization)
	r.GET("/organizations/:id", controllers.GetOrganization)
	r.DELETE("/organizations/:id", controllers.DeleteOrganization)
	r.PATCH("/organizations/:id", controllers.UpdateOrganization)

	r.Run()
}
