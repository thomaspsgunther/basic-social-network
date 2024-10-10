package comments

import (
	"time"
	"y_net/internal/services/shared"

	"github.com/google/uuid"
)

type Comments struct {
	CommentList []*Comment
}

type Comment struct {
	ID          uuid.UUID    `json:"id"`
	User        *shared.User `json:"user"`
	PostID      uuid.UUID    `json:"postId"`
	Description string       `json:"description"`
	LikeCount   int          `json:"likeCount"`
	CreatedAt   time.Time    `json:"createdAt"`
}
