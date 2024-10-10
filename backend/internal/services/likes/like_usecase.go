package likes

import (
	"context"
	"y_net/internal/services/shared"

	"github.com/google/uuid"
)

type LikeUsecase interface {
	LikePost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error
	UnlikePost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error
	UserLikedPost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error)
	GetFromPost(ctx context.Context, postId uuid.UUID) ([]shared.User, error)
}

type likeUsecaseImpl struct {
	usecase    LikeUsecase
	repository likeRepository
}

func NewLikeUsecase() LikeUsecase {
	return &likeUsecaseImpl{
		usecase:    &likeUsecaseImpl{},
		repository: &likeRepositoryImpl{},
	}
}

func (i *likeUsecaseImpl) LikePost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error {
	err := i.repository.likePost(ctx, userId, postId)
	if err != nil {
		return err
	}

	return nil
}

func (i *likeUsecaseImpl) UnlikePost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error {
	err := i.repository.unlikePost(ctx, userId, postId)
	if err != nil {
		return err
	}

	return nil
}

func (i *likeUsecaseImpl) UserLikedPost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error) {
	isLiked, err := i.repository.userLikedPost(ctx, userId, postId)
	if err != nil {
		return false, err
	}

	return isLiked, nil
}

func (i *likeUsecaseImpl) GetFromPost(ctx context.Context, postId uuid.UUID) ([]shared.User, error) {
	users, err := i.repository.getFromPost(ctx, postId)
	if err != nil {
		return nil, err
	}

	return users, nil
}
