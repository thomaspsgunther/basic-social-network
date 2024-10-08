package likes

import "github.com/google/uuid"

type Like struct {
	UserID uuid.UUID `json:"userId"`
	PostID uuid.UUID `json:"postId"`
}

type LikedJson struct {
	Liked bool `json:"liked"`
}
