package create_post

import (
	database "postservice/internal/db"
	objectstorage "postservice/internal/objectStorage"
)

type CreatePostRepository struct {
	dataRepository   *database.Database
	objectRepository *objectstorage.ObjectStorage
}

func NewCreatePostRepository(dataRepository *database.Database, objectRepository *objectstorage.ObjectStorage) *CreatePostRepository {
	return &CreatePostRepository{
		dataRepository:   dataRepository,
		objectRepository: objectRepository,
	}
}

type PostKey struct {
	PostId string
}

type PostMetadata struct {
	PostId       string `json:"post_id"`
	User         string `json:"username"`
	Type         string `json:"type"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Size         int    `json:"size"`
	HasThumbnail bool   `json:"has_thumbnail"`
	CreatedAt    string `json:"created_at"`
	LastUpdated  string `json:"last_updated"`
}

func (r *CreatePostRepository) AddNewPostMetaData(post *Post) error {
	data := &PostMetadata{
		PostId:       post.PostId,
		User:         post.User,
		Type:         post.Type,
		Title:        post.Title,
		Description:  post.Description,
		HasThumbnail: post.HasThumbnail,
		CreatedAt:    post.CreatedAt,
		LastUpdated:  post.LastUpdated,
	}
	return r.dataRepository.Client.InsertData("Posts", data)
}

func (r *CreatePostRepository) GetPresignedUrlsForUploading(post *Post) (PresignedUrl, error) {
	key := post.User + "/" + post.Type + "/" + post.PostId
	var presignedUrl PresignedUrl
	uploadId, contentPresignedUrls, err := r.objectRepository.Client.GetPreSignedUrlsForPuttingObject(key, post.Size)
	presignedUrl.UploadId = uploadId
	presignedUrl.ContentPresignedUrls = contentPresignedUrls

	if err != nil {
		return PresignedUrl{}, err
	}

	if post.HasThumbnail {
		thumbnailKey := post.User + "/" + post.Type + "/THUMBNAILS/" + post.PostId
		_, thumbanilPresignedUrl, err := r.objectRepository.Client.GetPreSignedUrlsForPuttingObject(thumbnailKey, 0)

		if err != nil {
			return PresignedUrl{}, err
		}

		presignedUrl.ThumbanilPresignedUrl = thumbanilPresignedUrl[0]
	}
	return presignedUrl, nil
}

func (r *CreatePostRepository) GetPostMetadata(postId string) (*Post, error) {
	postKey := &PostKey{
		PostId: postId,
	}
	var post Post
	err := r.dataRepository.Client.GetData("Posts", postKey, &post)

	return &post, err
}

func (r *CreatePostRepository) CompleteMultipartUpload(multipartPost *MultipartPost) error {
	multipartObject := convertMultipartPostToMultipartObject(multipartPost)
	return r.objectRepository.Client.CompleteMultipartUpload(multipartObject)
}

func (r *CreatePostRepository) RemoveUnconfirmedPost(postId string) error {
	postKey := &PostKey{
		PostId: postId,
	}
	return r.dataRepository.Client.RemoveData("Posts", postKey)
}

func convertMultipartPostToMultipartObject(post *MultipartPost) objectstorage.MultipartObject {
	key := post.Post.User + "/" + post.Post.Type + "/" + post.Post.PostId
	completedParts := make([]objectstorage.CompletedPart, len(post.CompletedParts))
	for i, part := range post.CompletedParts {
		completedParts[i] = objectstorage.CompletedPart{
			PartNumber: part.PartNumber,
			ETag:       part.ETag,
		}
	}

	return objectstorage.MultipartObject{
		Key:           key,
		UploadID:      post.UploadId,
		CompletedPart: completedParts,
	}
}
