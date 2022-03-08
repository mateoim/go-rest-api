package main

import (
	"github.com/gin-gonic/gin"
	"go-rest-api/models"
	"net/http"
)

func main() {
	r := gin.Default()
	models.ConnectToDatabase()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "hello world"})
	})

	r.Run()
}
