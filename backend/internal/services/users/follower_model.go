package users

import "github.com/google/uuid"

type Follower struct {
	FollowerID uuid.UUID `json:"followerId,omitempty"`
	FollowedID uuid.UUID `json:"followedId,omitempty"`
}

type FollowsJson struct {
	Follows bool `json:"follows,omitempty"`
}
