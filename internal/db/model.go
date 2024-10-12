package database

import "time"

type PostKey struct {
	PostId string
}

type Post struct {
	PostId       string    `json:"post_id"`
	User         string    `json:"username"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	LastUpdated  time.Time `json:"last_updated"`
	HasThumbnail bool      `json:"has_thumbnail"`
}
