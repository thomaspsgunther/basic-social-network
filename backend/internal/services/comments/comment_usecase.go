package comments

import (
	"context"
	"fmt"
	"y_net/internal/services/shared"

	"github.com/google/uuid"
)

type CommentUsecase interface {
	Create(ctx context.Context, comment Comment) (uuid.UUID, error)
	GetFromPost(ctx context.Context, postId uuid.UUID) ([]Comment, error)
	Get(ctx context.Context, id uuid.UUID) (Comment, error)
	Update(ctx context.Context, comment Comment, id uuid.UUID) error
	Like(ctx context.Context, id uuid.UUID) error
	Unlike(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type commentUsecaseImpl struct {
	usecase    CommentUsecase
	repository commentRepository
}

func NewCommentUsecase() CommentUsecase {
	return &commentUsecaseImpl{
		usecase:    &commentUsecaseImpl{},
		repository: &commentRepositoryImpl{},
	}
}

func (u *commentUsecaseImpl) Create(ctx context.Context, comment Comment) (uuid.UUID, error) {
	if (comment.User == &shared.User{}) {
		return uuid.Nil, fmt.Errorf("user must not be empty")
	}
	if comment.PostID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("postid must not be empty")
	}
	if comment.Description == "" {
		return uuid.Nil, fmt.Errorf("comment text must not be empty")
	}

	id, err := u.repository.create(ctx, comment)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (u *commentUsecaseImpl) GetFromPost(ctx context.Context, postId uuid.UUID) ([]Comment, error) {
	comments, err := u.repository.getFromPost(ctx, postId)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (u *commentUsecaseImpl) Get(ctx context.Context, id uuid.UUID) (Comment, error) {
	comment, err := u.repository.get(ctx, id)
	if err != nil {
		return Comment{}, err
	}

	return comment, nil
}

func (u *commentUsecaseImpl) Update(ctx context.Context, comment Comment, id uuid.UUID) error {
	err := u.repository.update(ctx, comment, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *commentUsecaseImpl) Like(ctx context.Context, id uuid.UUID) error {
	err := u.repository.like(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *commentUsecaseImpl) Unlike(ctx context.Context, id uuid.UUID) error {
	err := u.repository.unlike(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *commentUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.repository.delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
