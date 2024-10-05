package users

import "github.com/google/uuid"

type Users struct {
	UserList []*User
}

type User struct {
	ID            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	Password      string    `json:"password,omitempty"`
	Email         *string   `json:"email"`
	FullName      *string   `json:"fullName"`
	Description   *string   `json:"description"`
	Avatar        *string   `json:"avatar"`
	FollowerCount int       `json:"followerCount"`
}

type TokenJson struct {
	Token string `json:"token"`
}
