package posts

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestSetup struct {
	usecase PostUsecaseI
	repo    *mockPostRepository
}

func setup() *TestSetup {
	repo := newMockPostRepository()
	usecase := &postUsecase{repository: repo}

	return &TestSetup{usecase: usecase, repo: repo}
}
func TestCreatePost(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	post := Post{UserID: userId, Image: "image_url.jpg"}

	id, err := ts.usecase.Create(post)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestCreatePostEmptyImage(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	post := Post{UserID: userId, Image: ""}

	id, err := ts.usecase.Create(post)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestGetPost(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	post := Post{UserID: userId, Image: "image.jpg"}
	id, _ := ts.repo.create(post)

	retrievedPost, err := ts.usecase.GetPost(id)
	assert.NoError(t, err)
	assert.Equal(t, post.Image, retrievedPost.Image)
}

func TestGetPostNotFound(t *testing.T) {
	ts := setup()

	id := uuid.New()

	_, err := ts.usecase.GetPost(id)
	assert.Error(t, err)
}

func TestUpdatePost(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	post := Post{UserID: userId, Image: "image.jpg"}
	id, _ := ts.repo.create(post)

	post.Image = "updated_image.jpg"
	err := ts.usecase.Update(post, id)
	assert.NoError(t, err)

	updatedPost, err := ts.usecase.GetPost(id)
	assert.NoError(t, err)
	assert.Equal(t, "updated_image.jpg", updatedPost.Image)
}

func TestUpdatePostEmptyImage(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	post := Post{UserID: userId, Image: "image.jpg"}
	id, _ := ts.repo.create(post)

	post.Image = ""
	err := ts.usecase.Update(post, id)
	assert.Error(t, err)
}

func TestDeletePost(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	post := Post{UserID: userId, Image: "image.jpg"}
	id, _ := ts.repo.create(post)

	err := ts.usecase.Delete(id)
	assert.NoError(t, err)

	posts, err := ts.usecase.GetPosts(10, time.Now(), uuid.Nil)
	assert.NoError(t, err)
	assert.NotContains(t, posts, post)
}

func TestDeletePostNotFound(t *testing.T) {
	ts := setup()

	id := uuid.New()

	err := ts.usecase.Delete(id)
	assert.Error(t, err)
	assert.Equal(t, "post not found", err.Error())
}

// mockPostRepository is a mock implementation of postRepositoryI for testing
type mockPostRepository struct {
	posts map[uuid.UUID]Post
}

func newMockPostRepository() *mockPostRepository {
	return &mockPostRepository{
		posts: make(map[uuid.UUID]Post),
	}
}

func (m *mockPostRepository) create(post Post) (uuid.UUID, error) {
	id := uuid.New()
	m.posts[id] = post

	return id, nil
}

func (m *mockPostRepository) getPosts(limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]Post, error) {
	var result []Post
	for id, post := range m.posts {
		if len(result) >= limit {
			break
		}
		if id != lastId && post.CreatedAt.After(lastCreatedAt) {
			result = append(result, post)
		}
	}

	return result, nil
}

func (m *mockPostRepository) getPost(id uuid.UUID) (Post, error) {
	post, exists := m.posts[id]
	if !exists {
		return Post{}, fmt.Errorf("post not found")
	}

	return post, nil
}

func (m *mockPostRepository) update(post Post, id uuid.UUID) error {
	if _, exists := m.posts[id]; !exists {
		return fmt.Errorf("post not found")
	}
	m.posts[id] = post

	return nil
}

func (m *mockPostRepository) delete(id uuid.UUID) error {
	if _, exists := m.posts[id]; !exists {
		return fmt.Errorf("post not found")
	}
	delete(m.posts, id)

	return nil
}
func (m *mockPostRepository) getFromUser(userId uuid.UUID, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]Post, error) {
	var result []Post
	for _, post := range m.posts {
		if len(result) >= limit {
			break
		}
		if post.UserID == userId && post.CreatedAt.After(lastCreatedAt) && post.ID != lastId {
			result = append(result, post)
		}
	}

	return result, nil
}
