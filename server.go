package main

import (
	"io"
	"net/http"

	"github.com/Amitmahato/golang-file-upload-server/middleware"
	"github.com/Amitmahato/golang-file-upload-server/service"
	"github.com/gin-gonic/gin"
)

func main(){
	httpRouter := gin.Default()
	httpRouter.Use(middleware.CORSMiddleware())

	httpRouter.GET("/ping",func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{
			"message":"pong",
		})
	})

	httpRouter.POST("/upload",func(c *gin.Context) {
		file, fileHeader, err := c.Request.FormFile("file")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "No file is received",
			})
			return
		}

		// first 512 byte of file is supposed to contain file metadata like file header
		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
			})
			return
		}

		// detect the content-type of the file uploaded, allow only if its image file
		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "The provided file format is not allowed. Please upload a JPEG or PNG image",
			})
			return
		}

		// seek the file cursor to start of the file
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		// Upload the file to google cloud storage bucket
		storageDir := "images"	// images directory inside the storage bucket - directory should be appended with filename before uploading
		fileHeader.Filename = storageDir+"/"+fileHeader.Filename

		storageService := service.NewGCPStorageService()

		url,_ := storageService.UploadFile(c.Request,file,fileHeader)

		c.JSON(http.StatusOK,map[string]string{
			"url":url.String(),
		})
	})

	httpRouter.Run(":8000")
}