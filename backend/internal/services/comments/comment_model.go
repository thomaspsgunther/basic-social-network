package comments

import (
	"time"
	"y-net/internal/services/shared"

	"github.com/google/uuid"
)

type Comments struct {
	CommentList []*Comment
}

type Comment struct {
	ID        uuid.UUID    `json:"id,omitempty"`
	User      *shared.User `json:"user,omitempty"`
	PostID    uuid.UUID    `json:"postId,omitempty"`
	Message   string       `json:"message,omitempty"`
	CreatedAt time.Time    `json:"createdAt,omitempty"`
}
