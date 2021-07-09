package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main(){
	httpRouter := gin.Default()

	httpRouter.GET("/ping",func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"message":"pong",
		})
	})

	httpRouter.Run(":8000")
}