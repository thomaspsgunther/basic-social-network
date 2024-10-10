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

func (i *commentUsecaseImpl) Create(ctx context.Context, comment Comment) (uuid.UUID, error) {
	if (comment.User == &shared.User{}) {
		return uuid.Nil, fmt.Errorf("user must not be empty")
	}
	if comment.PostID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("postid must not be empty")
	}
	if comment.Description == "" {
		return uuid.Nil, fmt.Errorf("comment text must not be empty")
	}

	id, err := i.repository.create(ctx, comment)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (i *commentUsecaseImpl) GetFromPost(ctx context.Context, postId uuid.UUID) ([]Comment, error) {
	comments, err := i.repository.getFromPost(ctx, postId)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (i *commentUsecaseImpl) Get(ctx context.Context, id uuid.UUID) (Comment, error) {
	comment, err := i.repository.get(ctx, id)
	if err != nil {
		return Comment{}, err
	}

	return comment, nil
}

func (i *commentUsecaseImpl) Update(ctx context.Context, comment Comment, id uuid.UUID) error {
	err := i.repository.update(ctx, comment, id)
	if err != nil {
		return err
	}

	return nil
}

func (i *commentUsecaseImpl) Like(ctx context.Context, id uuid.UUID) error {
	err := i.repository.like(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (i *commentUsecaseImpl) Unlike(ctx context.Context, id uuid.UUID) error {
	err := i.repository.unlike(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (i *commentUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := i.repository.delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
