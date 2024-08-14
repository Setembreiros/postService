package objectstorage

//go:generate mockgen -source=object_storage.go -destination=mock/object_storage.go

type ObjectStorage struct {
	Client     ObjectStorageClient
	BucketName string
}

type ObjectStorageClient interface {
	GetPreSignedUrlForPuttingObject(objectKey string) (string, error)
}

func NewObjectStorage(client ObjectStorageClient) *ObjectStorage {
	return &ObjectStorage{
		Client: client,
	}
}