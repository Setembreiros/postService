package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

func (s3c *S3Client) GetPreSignedUrlForGettingObject(objectKey string) (string, error) {
	request, err := s3c.presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s3c.bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(s3c.presignLifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Couldn't get a presigned request to get %v:%v.",
			s3c.bucketName, objectKey)
	}
	return request.URL, err
}

func (s3c *S3Client) DeleteObjects(objectKeys []string) error {
	objects := make([]types.ObjectIdentifier, len(objectKeys))
	for i, key := range objectKeys {
		objects[i] = types.ObjectIdentifier{
			Key: aws.String(key),
		}
	}

	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(s3c.bucketName),
		Delete: &types.Delete{
			Objects: objects,
			Quiet:   aws.Bool(true), // Set to true to not receive a list of deleted objects
		},
	}

	_, err := s3c.client.DeleteObjects(context.TODO(), input)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to delete objects %v", objectKeys)
		return err
	}

	return nil
}
