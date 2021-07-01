package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/Amitmahato/file-server/utils"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
)

type gcpStorageService struct {
	StorageClient *storage.Client
}

// GCPStorageService Interface
type GCPStorageService interface {
	UploadFile(req *http.Request,file multipart.File, fileHeader *multipart.FileHeader) (*url.URL, error)
}

// NewFirebaseAuthService constructor
func NewGCPStorageService(sc *storage.Client) GCPStorageService {
	return &gcpStorageService{
		StorageClient: sc,
	}
}

const googleStorage = "https://storage.googleapis.com/"

func (sc *gcpStorageService) UploadFile(req *http.Request,file multipart.File, fileHeader *multipart.FileHeader) (*url.URL, error) {
	fileName := fileHeader.Filename

	directory := strings.Split(fileName,"/")
	nameAndExtension := strings.Split(directory[len(directory)-1],".")
	fileExtension := nameAndExtension[len(nameAndExtension)-1]
	uniqueId := uuid.New()

	fileName = strings.Replace(uniqueId.String(), "-", "", -1)
	fileName = fmt.Sprintf("%v.%v",fileName,fileExtension)
	directory[len(directory)-1] = fileName

	fileName = strings.Join(directory,"/")

	bucketName := utils.GetEnvWithKey("STORAGE_BUCKET_NAME")
	bucket := sc.StorageClient.Bucket(bucketName)
	ctx := req.Context()
	obj := bucket.Object(fileName).NewWriter(ctx)
	// obj.Metadata = map[string]string{
	// 	"x-goog-acl":"public-read",
	// }

	obj.ContentType = fileHeader.Header.Get("Content-Type")
	if _, err := io.Copy(obj, file); err != nil {
		return nil,err
	}
	acl := bucket.Object(obj.Name).ACL()
	if err := obj.Close(); err != nil {
		return nil,err
	}
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return nil, err
	}
	return url.Parse(googleStorage + bucketName + "/" + obj.Attrs().Name)
}
