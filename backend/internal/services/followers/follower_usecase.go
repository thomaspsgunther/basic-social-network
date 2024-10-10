package followers

import (
	"context"
	"y_net/internal/services/shared"

	"github.com/google/uuid"
)

type FollowerUsecase interface {
	Follow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error
	Unfollow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error
	UserFollowsUser(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) (bool, error)
	GetFollowers(ctx context.Context, userId uuid.UUID) ([]shared.User, error)
	GetFollowed(ctx context.Context, userId uuid.UUID) ([]shared.User, error)
}

type followerUsecaseImpl struct {
	usecase    FollowerUsecase
	repository followerRepository
}

func NewFollowerUsecase() FollowerUsecase {
	return &followerUsecaseImpl{
		usecase:    &followerUsecaseImpl{},
		repository: &followerRepositoryImpl{},
	}
}

func (u *followerUsecaseImpl) Follow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error {
	err := u.repository.follow(ctx, followerId, followedId)
	if err != nil {
		return err
	}

	return nil
}

func (u *followerUsecaseImpl) Unfollow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error {
	err := u.repository.unfollow(ctx, followerId, followedId)
	if err != nil {
		return err
	}

	return nil
}

func (u *followerUsecaseImpl) UserFollowsUser(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) (bool, error) {
	follows, err := u.repository.userFollowsUser(ctx, followerId, followedId)
	if err != nil {
		return false, err
	}

	return follows, nil
}

func (u *followerUsecaseImpl) GetFollowers(ctx context.Context, userId uuid.UUID) ([]shared.User, error) {
	users, err := u.repository.getFollowers(ctx, userId)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *followerUsecaseImpl) GetFollowed(ctx context.Context, userId uuid.UUID) ([]shared.User, error) {
	users, err := u.repository.getFollowed(ctx, userId)
	if err != nil {
		return nil, err
	}

	return users, nil
}
