package routes

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

type Routes struct {
	handler infrastructures.GinRouter
	storageService services.GCPStorageService
}

func NewRoutes(handler infrastructures.GinRouter, storageService services.GCPStorageService) Routes{
	return Routes{
		handler: handler,
		storageService: storageService,
	}

	
}

func (r *Routes) Setup(){
	
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.handler.Gin.MaxMultipartMemory = 8 << 20  // 8 MiB
	r.handler.Gin.POST("/upload", func(c *gin.Context) {
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
		url,_ := r.storageService.UploadFile(c.Request,file,fileHeader)

		c.JSON(http.StatusOK,map[string]string{
			"url":url.String(),
		})
	})
}