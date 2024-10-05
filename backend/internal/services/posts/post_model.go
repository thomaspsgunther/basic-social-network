package posts

import (
	"time"

	"github.com/google/uuid"
)

type Posts struct {
	PostList []*Post
}

type Post struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userId"`
	Image        string    `json:"image"`
	Description  *string   `json:"description"`
	LikeCount    int       `json:"likeCount"`
	CommentCount int       `json:"commentCount"`
	CreatedAt    time.Time `json:"createdAt"`
}
