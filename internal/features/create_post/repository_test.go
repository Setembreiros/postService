package create_post_test

import (
	database "postservice/internal/db"
	mock_database "postservice/internal/db/mock"
	"postservice/internal/features/create_post"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

var client *mock_database.MockDatabaseClient
var createPostRepository create_post.CreatePostRepository

func setUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	client = mock_database.NewMockDatabaseClient(ctrl)
	createPostRepository = create_post.CreatePostRepository(*database.NewDatabase(client))
}

func TestAddNewPostMetaDataInRepository(t *testing.T) {
	setUp(t)
	newPost := &create_post.Post{
		User:        "username1",
		Title:       "Meu Post",
		Description: "Este é o meu novo post",
		CreatedAt:   time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
		LastUpdated: time.Date(2024, 8, 8, 21, 51, 20, 33, time.UTC).UTC(),
	}
	postId := "username1-Meu_Post-1723153880"
	data := &create_post.PostMetadata{
		PostId:   postId,
		Metadata: newPost,
	}
	client.EXPECT().InsertData("Posts", data)

	createPostRepository.AddNewPostMetaData(postId, newPost)
}
