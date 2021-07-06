package infrastructures

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Middleware1() gin.HandlerFunc {
	log.Println("Demo Middleware1 Registered")
	return func(c *gin.Context) {
		log.Println("Demo Middleware1 Executed")
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	log.Println("CORS Middleware2 Registered")
	return func(c *gin.Context) {
		// CORS middleware
		// this will allow direct usage of axios requests from next js app
		// no more request forwarding from api/ directory a.k.a request forwarding for API (serverIP/*) through API routes (/api/*)
		log.Println("CORS Middleware2 Executed")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3001")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func Middleware2() gin.HandlerFunc {
	log.Println("Demo Middleware3 Registered")
	return func(c *gin.Context) {
		log.Println("Demo Middleware3 Executed")
		c.Next()
	}
}

type GinRouter struct {
	Gin *gin.Engine
}

func NewRouter() GinRouter {
	router := gin.Default()
	
	router.Use(Middleware1())
	router.Use(CORSMiddleware())
	router.Use(Middleware2())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	
	return GinRouter{
		Gin: router,
	}
}