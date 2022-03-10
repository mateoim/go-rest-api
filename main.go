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

	r.GET("/users", controllers.GetUsers)
	r.POST("/users", controllers.CreateUser)
	r.GET("/users/:id", controllers.GetUser)
	r.DELETE("/users/:id", controllers.DeleteUser)
	r.PATCH("/users/:id", controllers.UpdateUser)

	r.GET("/organizations/:id/users", controllers.GetUsersByOrganization)

	r.Run()
}
