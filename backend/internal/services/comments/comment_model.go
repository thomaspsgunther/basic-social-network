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
	ID        uuid.UUID    `json:"id"`
	User      *shared.User `json:"user"`
	PostID    uuid.UUID    `json:"postId"`
	Message   string       `json:"message"`
	LikeCount int          `json:"likeCount"`
	CreatedAt time.Time    `json:"createdAt"`
}
