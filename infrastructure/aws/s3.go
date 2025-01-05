package aws

import (
	"context"
	"math"
	objectstorage "postservice/internal/objectStorage"
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

func (s3c *S3Client) GetPreSignedUrlsForPuttingObject(objectKey string, size int) ([]string, error) {
	if size > 100 {
		return s3c.getMultipartPreSignedUrls(objectKey, size)
	}

	presignedUrl, err := s3c.getPreSignedUrl(objectKey)
	return []string{presignedUrl}, err
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

func (s3c *S3Client) CompleteMultipartUpload(multipartObject objectstorage.MultipartObject) error {
	completeMultipartInput := &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(s3c.bucketName),
		Key:      aws.String(multipartObject.Key),
		UploadId: aws.String(multipartObject.UploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: transformCompletedParts(multipartObject.CompletedPart),
		},
	}

	_, err := s3c.client.CompleteMultipartUpload(context.TODO(), completeMultipartInput)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to complete multipart upload")
	}

	return err
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

func (s3c *S3Client) getPreSignedUrl(objectKey string) (string, error) {
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

func (s3c *S3Client) getMultipartPreSignedUrls(objectKey string, size int) ([]string, error) {
	createMultipartUploadInput := &s3.CreateMultipartUploadInput{
		Bucket: aws.String(s3c.bucketName),
		Key:    aws.String(objectKey),
	}

	multipartOutput, err := s3c.client.CreateMultipartUpload(context.TODO(), createMultipartUploadInput)
	if err != nil {
		log.Error().Stack().Err(err).Msgf("Failed to initiate multipart upload")
		return []string{}, err
	}

	uploadID := *multipartOutput.UploadId
	log.Info().Msgf("Multipart upload iniciado. UploadID: %s\n", uploadID)

	numParts := int(math.Ceil(float64(size) / 100))
	presinedUrls := []string{}

	for part := 1; part <= numParts; part++ {
		request, err := s3c.presignClient.PresignUploadPart(context.TODO(), &s3.UploadPartInput{
			Bucket:     aws.String(s3c.bucketName),
			Key:        aws.String(objectKey),
			PartNumber: aws.Int32(int32(part)),
			UploadId:   aws.String(uploadID),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(s3c.presignLifetimeSecs * int64(time.Second))
		})
		if err != nil {
			log.Error().Stack().Err(err).Msgf("Couldn't get a presigned request to put %v:%v.",
				s3c.bucketName, objectKey)
			return []string{}, err
		}

		presinedUrls = append(presinedUrls, request.URL)
	}

	return presinedUrls, err
}

func transformCompletedParts(parts []objectstorage.CompletedPart) []types.CompletedPart {
	// Slice para almacenar o resultado
	var s3Parts []types.CompletedPart

	for _, part := range parts {
		s3Parts = append(s3Parts, types.CompletedPart{
			PartNumber: aws.Int32(int32(part.PartNumber)),
			ETag:       aws.String(part.ETag),
		})
	}

	return s3Parts
}
