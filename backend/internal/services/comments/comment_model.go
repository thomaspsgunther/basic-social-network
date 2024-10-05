package comments

import (
	"time"

	"github.com/google/uuid"
)

type Comments struct {
	CommentList []*Comment
}

type Comment struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	PostID      uuid.UUID `json:"postId"`
	Description string    `json:"description"`
	LikeCount   int       `json:"likeCount"`
	CreatedAt   time.Time `json:"createdAt"`
}
