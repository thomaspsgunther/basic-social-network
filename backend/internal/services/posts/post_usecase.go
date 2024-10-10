package posts

import (
	"context"
	"fmt"
	"time"
	"y_net/internal/services/shared"

	"github.com/google/uuid"
)

type IPostUsecase interface {
	Create(ctx context.Context, post shared.Post) (uuid.UUID, error)
	GetPosts(ctx context.Context, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]shared.Post, error)
	GetPost(ctx context.Context, id uuid.UUID) (shared.Post, error)
	Update(ctx context.Context, post shared.Post, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	Like(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error
	Unlike(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error
	UserLikedPost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error)
	GetLikes(ctx context.Context, id uuid.UUID) ([]shared.User, error)
}

type postUsecaseImpl struct {
	usecase    IPostUsecase
	repository iPostRepository
}

func NewPostUsecase() IPostUsecase {
	return &postUsecaseImpl{
		usecase:    &postUsecaseImpl{},
		repository: &postRepositoryImpl{},
	}
}

func (u *postUsecaseImpl) Create(ctx context.Context, post shared.Post) (uuid.UUID, error) {
	if (post.User == &shared.User{}) {
		return uuid.Nil, fmt.Errorf("user must not be empty")
	}
	if post.Image == "" {
		return uuid.Nil, fmt.Errorf("post image must not be empty")
	}

	id, err := u.repository.create(ctx, post)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (u *postUsecaseImpl) GetPosts(ctx context.Context, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]shared.Post, error) {
	posts, err := u.repository.getPosts(ctx, limit, lastCreatedAt, lastId)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (u *postUsecaseImpl) GetPost(ctx context.Context, id uuid.UUID) (shared.Post, error) {
	post, err := u.repository.getPost(ctx, id)
	if err != nil {
		return shared.Post{}, err
	}

	return post, nil
}

func (u *postUsecaseImpl) Update(ctx context.Context, post shared.Post, id uuid.UUID) error {
	if post.Image == "" {
		return fmt.Errorf("post image must not be empty")
	}

	err := u.repository.update(ctx, post, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *postUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.repository.delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (i *postUsecaseImpl) Like(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error {
	err := i.repository.like(ctx, userId, postId)
	if err != nil {
		return err
	}

	return nil
}

func (i *postUsecaseImpl) Unlike(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error {
	err := i.repository.unlike(ctx, userId, postId)
	if err != nil {
		return err
	}

	return nil
}

func (i *postUsecaseImpl) UserLikedPost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error) {
	isLiked, err := i.repository.userLikedPost(ctx, userId, postId)
	if err != nil {
		return false, err
	}

	return isLiked, nil
}

func (i *postUsecaseImpl) GetLikes(ctx context.Context, id uuid.UUID) ([]shared.User, error) {
	users, err := i.repository.getLikes(ctx, id)
	if err != nil {
		return nil, err
	}

	return users, nil
}
