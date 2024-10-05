package posts

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PostUsecaseI interface {
	Create(post Post) (uuid.UUID, error)
	GetPosts(limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]Post, error)
	GetPost(id uuid.UUID) (Post, error)
	Update(post Post, id uuid.UUID) error
	Delete(id uuid.UUID) error
	GetFromUser(userId uuid.UUID, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]Post, error)
}

type postUsecase struct {
	usecase    PostUsecaseI
	repository postRepositoryI
}

func NewPostUsecase() PostUsecaseI {
	return &postUsecase{
		usecase:    &postUsecase{},
		repository: &postRepository{},
	}
}

func (i *postUsecase) Create(post Post) (uuid.UUID, error) {
	if post.Image == "" {
		return uuid.UUID{}, fmt.Errorf("post image must not be empty")
	}

	id, err := i.repository.create(post)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (i *postUsecase) GetPosts(limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]Post, error) {
	posts, err := i.repository.getPosts(limit, lastCreatedAt, lastId)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (i *postUsecase) GetPost(id uuid.UUID) (Post, error) {
	post, err := i.repository.getPost(id)
	if err != nil {
		return Post{}, err
	}

	return post, nil
}

func (i *postUsecase) Update(post Post, id uuid.UUID) error {
	if post.Image == "" {
		return fmt.Errorf("post image must not be empty")
	}

	err := i.repository.update(post, id)
	if err != nil {
		return err
	}

	return nil
}

func (i *postUsecase) Delete(id uuid.UUID) error {
	err := i.repository.delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (i *postUsecase) GetFromUser(userId uuid.UUID, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]Post, error) {
	posts, err := i.repository.getFromUser(userId, limit, lastCreatedAt, lastId)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
