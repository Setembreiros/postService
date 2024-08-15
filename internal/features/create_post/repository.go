package create_post

import (
	"fmt"
	database "postservice/internal/db"
	objectstorage "postservice/internal/objectStorage"
	"strconv"
	"strings"
	"time"
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
	PostId      string    `json:"post_id"`
	User        string    `json:"username"`
	Type        string    `json:"type"`
	FileType    string    `json:"file_type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	LastUpdated time.Time `json:"last_updated"`
}

func (r *CreatePostRepository) AddNewPostMetaData(post *Post) error {
	data := &PostMetadata{
		PostId:      generatePostId(post),
		User:        post.User,
		Type:        post.Type,
		FileType:    post.FileType,
		Title:       post.Title,
		Description: post.Description,
		CreatedAt:   post.CreatedAt,
		LastUpdated: post.LastUpdated,
	}
	return r.dataRepository.Client.InsertData("Posts", data)
}

func (r *CreatePostRepository) GetPresignedUrlForUploadingText(post *Post) (string, error) {
	key := post.User + "/" + post.Type + "/" + generatePostId(post) + "." + post.FileType
	return r.objectRepository.Client.GetPreSignedUrlForPuttingObject(key)
}

func (r *CreatePostRepository) GetPostMetadata(postId string) (*Post, error) {
	postKey := &PostKey{
		PostId: postId,
	}
	var post Post
	err := r.dataRepository.Client.GetData("Posts", postKey, &post)

	fmt.Println(post)
	return &post, err
}

func (r *CreatePostRepository) RemoveUnconfirmedPost(postId string) error {
	postKey := &PostKey{
		PostId: postId,
	}
	return r.dataRepository.Client.RemoveData("Posts", postKey)
}

func generatePostId(post *Post) string {
	postId := post.User + "-" + post.Title + "-" + strconv.FormatInt(post.CreatedAt.Unix(), 10)
	return strings.ReplaceAll(strings.ReplaceAll(postId, " ", "_"), "\t", "_")
}
