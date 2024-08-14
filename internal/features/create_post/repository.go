package create_post

import (
	database "postservice/internal/db"
	objectstorage "postservice/internal/objectStorage"
	"strconv"
	"strings"
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
	PostId   string `json:"post_id"`
	Metadata *Post  `json:"metadata"`
}

func (r CreatePostRepository) AddNewPostMetaData(post *Post) error {
	data := &PostMetadata{
		PostId:   generatePostId(post),
		Metadata: post,
	}
	return r.dataRepository.Client.InsertData("Posts", data)
}

func (r CreatePostRepository) GetPresignedUrlForUploadingText(post *Post) (string, error) {
	key := post.User + "/" + post.Type + "/" + generatePostId(post)
	return r.objectRepository.Client.GetPreSignedUrlForPuttingObject(key)
}

func (r CreatePostRepository) GetPostMetadata(postId string) (*Post, error) {
	postKey := &PostKey{
		PostId: postId,
	}
	var post Post
	err := r.dataRepository.Client.GetData("Posts", postKey, &post)

	return &post, err
}

func (r CreatePostRepository) RemoveUnconfirmedPost(postId string) error {
	return r.dataRepository.Client.RemoveData("Posts", postId)
}

func generatePostId(post *Post) string {
	postId := post.User + "-" + post.Title + "-" + strconv.FormatInt(post.CreatedAt.Unix(), 10)
	return strings.ReplaceAll(strings.ReplaceAll(postId, " ", "_"), "\t", "_")
}
