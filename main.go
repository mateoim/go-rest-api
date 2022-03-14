package main

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-rest-api/config"
	"go-rest-api/controllers"
	"go-rest-api/docs"
)

func main() {
	docs.SwaggerInfo.Title = "Swagger Go REST API"
	docs.SwaggerInfo.Description = "Simple Go REST API application."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r := SetupRouter()

	r.Run()
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	config.ConnectToDatabase()

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
	r.GET("/users/:id/events", controllers.GetUserEvents)
	r.GET("/users/:id/meetings", controllers.GetUserMeetings)
	r.GET("/users/:id/invitations", controllers.GetUserInvitations)

	r.GET("/organizations/:id/users", controllers.GetUsersByOrganization)

	r.GET("/events", controllers.GetEvents)
	r.POST("/events", controllers.CreateEvent)
	r.GET("/events/:id", controllers.GetEvent)
	r.DELETE("/events/:id", controllers.DeleteEvent)
	r.PATCH("/events/:id", controllers.UpdateEvent)
	r.POST("/events/:id/register", controllers.RegisterUser)
	r.GET("/events/:id/users", controllers.GetEventUsers)

	r.GET("/events/:id/meetings", controllers.GetMeetings)
	r.POST("/events/:id/meetings", controllers.CreateMeeting)
	r.GET("/events/:id/meetings/:meeting-id", controllers.GetMeeting)
	r.DELETE("/events/:id/meetings/:meeting-id", controllers.DeleteMeeting)
	r.POST("/events/:id/meetings/:meeting-id/schedule", controllers.Schedule)

	r.GET("/events/:id/meetings/:meeting-id/invitations", controllers.GetInvitations)
	r.POST("/events/:id/meetings/:meeting-id/invitations", controllers.CreateInvitation)
	r.GET("/events/:id/meetings/:meeting-id/invitations/:invitation-id", controllers.GetInvitation)
	r.DELETE("/events/:id/meetings/:meeting-id/invitations/:invitation-id", controllers.DeleteInvitation)
	r.POST("/events/:id/meetings/:meeting-id/invitations/:invitation-id/accept", controllers.Accept)
	r.POST("/events/:id/meetings/:meeting-id/invitations/:invitation-id/reject", controllers.Reject)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
