package main

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/config"
	"go-rest-api/controllers"
)

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

	r.GET("/organizations/:id/users", controllers.GetUsersByOrganization)

	r.GET("/events", controllers.GetEvents)
	r.POST("/events", controllers.CreateEvent)
	r.GET("/events/:id", controllers.GetEvent)
	r.DELETE("/events/:id", controllers.DeleteEvent)
	r.PATCH("/events/:id", controllers.UpdateEvent)
	r.POST("/events/:id/register", controllers.RegisterUser)

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

	return r
}

func main() {
	r := SetupRouter()

	r.Run()
}
