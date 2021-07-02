package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Amitmahato/file-server/infrastructures"
	"github.com/Amitmahato/file-server/services"
	"github.com/Amitmahato/file-server/utils"
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


func main() {
	router := gin.Default()

	utils.LoadEnv()

	storageClient := infrastructures.InitStorageClient()
	storageService := services.NewGCPStorageService(storageClient)

	router.Use(Middleware1())
	router.Use(CORSMiddleware())
	router.Use(Middleware2())

	
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// single file
		file, fileHeader, err := c.Request.FormFile("file")
		if err != nil {
			log.Println("ERROR : ",err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "No file is received",
			})
			return
		}

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
			})
			return
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "The provided file format is not allowed. Please upload a JPEG or PNG image",
			})
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
			})
			return
		}
		
		// Save the file locally to specific destination on server
		localDir := "images"
		os.MkdirAll(localDir,os.ModePerm)
		c.SaveUploadedFile(fileHeader, localDir+"/"+fileHeader.Filename)

		// Upload the file to google cloud storage bucket
		storageDir := utils.GetEnvWithKey("UPLOADED_FILE_DIR")
		fileHeader.Filename = storageDir+"/"+fileHeader.Filename
		url,_ := storageService.UploadFile(c.Request,file,fileHeader)

		c.JSON(http.StatusOK,map[string]string{
			"url":url.String(),
		})
	})
	router.Run(":8000")
}
