package shared

import "github.com/google/uuid"

type Users struct {
	UserList []*User
}

type User struct {
	ID            uuid.UUID `json:"id,omitempty"`
	Username      string    `json:"username,omitempty"`
	Password      string    `json:"password,omitempty"`
	Email         *string   `json:"email,omitempty"`
	FullName      *string   `json:"fullName,omitempty"`
	Description   *string   `json:"description,omitempty"`
	Avatar        *string   `json:"avatar,omitempty"`
	PostCount     int       `json:"postCount,omitempty"`
	FollowerCount int       `json:"followerCount,omitempty"`
	FollowedCount int       `json:"followedCount,omitempty"`
}
