package create_post

import database "postservice/internal/db"

type CreatePostRepository database.Database

type PostMetadata struct {
	PostId   string `json:"post_id"`
	Metadata *Post  `json:"metadata"`
}

func (r CreatePostRepository) AddNewPostMetaData(id string, post *Post) error {
	data := &PostMetadata{
		PostId:   id,
		Metadata: post,
	}
	return r.Client.InsertData("Posts", data)
}
