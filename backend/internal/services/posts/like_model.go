package posts

import "github.com/google/uuid"

type Like struct {
	UserID uuid.UUID `json:"userId,omitempty"`
	PostID uuid.UUID `json:"postId,omitempty"`
}

type LikedJson struct {
	Liked bool `json:"liked,omitempty"`
}
