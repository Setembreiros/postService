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

func (r *CreatePostRepository) GetPresignedUrlsForUploading(post *Post) ([]string, error) {
	urls := []string{}

	key := post.User + "/" + post.Type + "/" + post.PostId
	url, err := r.objectRepository.Client.GetPreSignedUrlForPuttingObject(key)
	if err != nil {
		return []string{}, err
	}
	urls = append(urls, url)
	if post.HasThumbnail {
		thumbnailKey := post.User + "/" + post.Type + "/THUMBNAILS/" + post.PostId
		thumbnailUrl, err := r.objectRepository.Client.GetPreSignedUrlForPuttingObject(thumbnailKey)
		if err != nil {
			return []string{}, err
		}
		urls = append(urls, thumbnailUrl)
	}
	return urls, nil
}

func (r *CreatePostRepository) GetPostMetadata(postId string) (*Post, error) {
	postKey := &PostKey{
		PostId: postId,
	}
	var post Post
	err := r.dataRepository.Client.GetData("Posts", postKey, &post)

	return &post, err
}

func (r *CreatePostRepository) RemoveUnconfirmedPost(postId string) error {
	postKey := &PostKey{
		PostId: postId,
	}
	return r.dataRepository.Client.RemoveData("Posts", postKey)
}
