package users

import "github.com/google/uuid"

type Follower struct {
	FollowerID uuid.UUID `json:"followerId"`
	FollowedID uuid.UUID `json:"followedId"`
}

type FollowsJson struct {
	Follows bool `json:"follows"`
}
