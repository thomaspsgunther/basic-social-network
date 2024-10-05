package comments

import (
	"fmt"

	"github.com/google/uuid"
)

type CommentUsecaseI interface {
	Create(comment Comment) (uuid.UUID, error)
	Update(comment Comment, id uuid.UUID) error
	Like(id uuid.UUID) error
	Unlike(id uuid.UUID) error
	Delete(id uuid.UUID) error
	GetFromPost(postId uuid.UUID) ([]Comment, error)
	Get(id uuid.UUID) (Comment, error)
}

type commentUsecase struct {
	usecase    CommentUsecaseI
	repository commentRepositoryI
}

func NewCommentUsecase() CommentUsecaseI {
	return &commentUsecase{
		usecase:    &commentUsecase{},
		repository: &commentRepository{},
	}
}

func (i *commentUsecase) Create(comment Comment) (uuid.UUID, error) {
	if comment.Description == "" {
		return uuid.UUID{}, fmt.Errorf("comment text must not be empty")
	}

	id, err := i.repository.create(comment)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (i *commentUsecase) Update(comment Comment, id uuid.UUID) error {
	err := i.repository.update(comment, id)
	if err != nil {
		return err
	}

	return nil
}

func (i *commentUsecase) Like(id uuid.UUID) error {
	err := i.repository.like(id)
	if err != nil {
		return err
	}

	return nil
}

func (i *commentUsecase) Unlike(id uuid.UUID) error {
	err := i.repository.unlike(id)
	if err != nil {
		return err
	}

	return nil
}

func (i *commentUsecase) Delete(id uuid.UUID) error {
	err := i.repository.delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (i *commentUsecase) GetFromPost(postId uuid.UUID) ([]Comment, error) {
	comments, err := i.repository.getFromPost(postId)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (i *commentUsecase) Get(id uuid.UUID) (Comment, error) {
	comment, err := i.repository.get(id)
	if err != nil {
		return Comment{}, err
	}

	return comment, nil
}
