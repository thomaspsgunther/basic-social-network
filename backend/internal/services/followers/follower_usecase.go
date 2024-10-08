package followers

import (
	"y_net/internal/services/users"

	"github.com/google/uuid"
)

type FollowerUsecaseI interface {
	Follow(followerId uuid.UUID, followedId uuid.UUID) error
	Unfollow(followerId uuid.UUID, followedId uuid.UUID) error
	UserFollowsUser(followerId uuid.UUID, followedId uuid.UUID) (bool, error)
	GetFollowers(userId uuid.UUID) ([]users.User, error)
	GetFollowed(userId uuid.UUID) ([]users.User, error)
}

type followerUsecase struct {
	usecase    FollowerUsecaseI
	repository followerRepositoryI
}

func NewFollowerUsecase() FollowerUsecaseI {
	return &followerUsecase{
		usecase:    &followerUsecase{},
		repository: &followerRepository{},
	}
}

func (i *followerUsecase) Follow(followerId uuid.UUID, followedId uuid.UUID) error {
	err := i.repository.follow(followerId, followedId)
	if err != nil {
		return err
	}

	return nil
}

func (i *followerUsecase) Unfollow(followerId uuid.UUID, followedId uuid.UUID) error {
	err := i.repository.unfollow(followerId, followedId)
	if err != nil {
		return err
	}

	return nil
}

func (i *followerUsecase) UserFollowsUser(followerId uuid.UUID, followedId uuid.UUID) (bool, error) {
	follows, err := i.repository.userFollowsUser(followerId, followedId)
	if err != nil {
		return false, err
	}

	return follows, nil
}

func (i *followerUsecase) GetFollowers(userId uuid.UUID) ([]users.User, error) {
	users, err := i.repository.getFollowers(userId)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (i *followerUsecase) GetFollowed(userId uuid.UUID) ([]users.User, error) {
	users, err := i.repository.getFollowed(userId)
	if err != nil {
		return nil, err
	}

	return users, nil
}
