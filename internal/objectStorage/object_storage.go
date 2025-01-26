package objectstorage

//go:generate mockgen -source=object_storage.go -destination=mock/object_storage.go

type ObjectStorage struct {
	Client     ObjectStorageClient
	BucketName string
}

type MultipartObject struct {
	Key           string          `json:"key"`
	UploadID      string          `json:"uploadID"`
	CompletedPart []CompletedPart `json:"completedPart"`
}

type CompletedPart struct {
	PartNumber int    `json:"partNumber"`
	ETag       string `json:"eTag"`
}

type ObjectStorageClient interface {
	GetPreSignedUrlsForPuttingObject(objectKey string, size int) (string, []string, error)
	GetPreSignedUrlForGettingObject(objectKey string) (string, error)
	CompleteMultipartUpload(multipartobject MultipartObject) error
	DeleteObjects(objectKeys []string) error
}

func NewObjectStorage(client ObjectStorageClient) *ObjectStorage {
	return &ObjectStorage{
		Client: client,
	}
}
