package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

type S3Client struct {
	client              *s3.Client
	presignClient       *s3.PresignClient
	presignLifetimeSecs int64
	bucketName          string
}

func NewS3Client(config aws.Config, bucketName string) *S3Client {
	s3Client := s3.NewFromConfig(config)
	return &S3Client{
		client:              s3Client,
		presignClient:       s3.NewPresignClient(s3Client),
		presignLifetimeSecs: 60,
		bucketName:          bucketName,
	}
}

func (s3c *S3Client) GetPreSignedUrlForPuttingObject(objectKey string) (string, error) {
	request, err := s3c.presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s3c.bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(s3c.presignLifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Couldn't get a presigned request to put %v:%v.",
			s3c.bucketName, objectKey)
	}
	return request.URL, err
}
