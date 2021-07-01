package infrastructures

import (
	"context"
	"path/filepath"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// InitStorageClient init the storage client
func InitStorageClient() *storage.Client {
	ctx := context.Background()
	serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")
	if err != nil {
		panic("Unable to load serviceAccountKey.json file")
	}
	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)
	client, err := storage.NewClient(ctx, opt)
	if err != nil {
		panic("Error creating GCP stroage client")
	}
	return client
}
