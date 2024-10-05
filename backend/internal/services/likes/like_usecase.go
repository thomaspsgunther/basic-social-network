package likes

import (
	"y_net/internal/services/users"

	"github.com/google/uuid"
)

type LikeUsecaseI interface {
	LikePost(userId uuid.UUID, postId uuid.UUID) error
	UnlikePost(userId uuid.UUID, postId uuid.UUID) error
	GetFromPost(postId uuid.UUID) ([]users.User, error)
}

type likeUsecase struct {
	usecase    LikeUsecaseI
	repository likeRepositoryI
}

func NewLikeUsecase() LikeUsecaseI {
	return &likeUsecase{
		usecase:    &likeUsecase{},
		repository: &likeRepository{},
	}
}

func (i *likeUsecase) LikePost(userId uuid.UUID, postId uuid.UUID) error {
	err := i.repository.likePost(userId, postId)
	if err != nil {
		return err
	}

	return nil
}

func (i *likeUsecase) UnlikePost(userId uuid.UUID, postId uuid.UUID) error {
	err := i.repository.unlikePost(userId, postId)
	if err != nil {
		return err
	}

	return nil
}

func (i *likeUsecase) GetFromPost(postId uuid.UUID) ([]users.User, error) {
	users, err := i.repository.getFromPost(postId)
	if err != nil {
		return nil, err
	}

	return users, nil
}
