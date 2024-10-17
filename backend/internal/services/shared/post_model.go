package shared

import (
	"time"

	"github.com/google/uuid"
)

type Posts struct {
	PostList []*Post
}

type Post struct {
	ID           uuid.UUID `json:"id,omitempty"`
	User         *User     `json:"user,omitempty"`
	Image        string    `json:"image,omitempty"`
	Description  *string   `json:"description,omitempty"`
	LikeCount    int       `json:"likeCount,omitempty"`
	CommentCount int       `json:"commentCount,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
}
