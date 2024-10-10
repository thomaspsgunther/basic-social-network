package posts

import (
	"context"
	"fmt"
	"testing"
	"time"
	"y_net/internal/services/shared"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestSetup struct {
	usecase PostUsecase
	repo    *mockPostRepository
}

func setup() *TestSetup {
	repo := newMockPostRepository()
	usecase := &postUsecaseImpl{repository: repo}

	return &TestSetup{usecase: usecase, repo: repo}
}
func TestCreatePost(t *testing.T) {
	ts := setup()

	user := shared.User{ID: uuid.New(), Username: "testuser"}
	post := shared.Post{User: &user, Image: "image_url.jpg"}

	id, err := ts.usecase.Create(context.Background(), post)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestCreatePostEmptyUser(t *testing.T) {
	ts := setup()

	post := shared.Post{User: &shared.User{}, Image: "image_url.jpg"}

	id, err := ts.usecase.Create(context.Background(), post)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestCreatePostEmptyImage(t *testing.T) {
	ts := setup()

	user := shared.User{ID: uuid.New(), Username: "testuser"}
	post := shared.Post{User: &user, Image: ""}

	id, err := ts.usecase.Create(context.Background(), post)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestGetPost(t *testing.T) {
	ts := setup()

	user := shared.User{ID: uuid.New(), Username: "testuser"}
	post := shared.Post{User: &user, Image: "image.jpg"}
	id, _ := ts.usecase.Create(context.Background(), post)

	retrievedPost, err := ts.usecase.GetPost(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, post.Image, retrievedPost.Image)
}

func TestGetPostNotFound(t *testing.T) {
	ts := setup()

	id := uuid.New()

	_, err := ts.usecase.GetPost(context.Background(), id)
	assert.Error(t, err)
}

func TestUpdatePost(t *testing.T) {
	ts := setup()

	user := shared.User{ID: uuid.New(), Username: "testuser"}
	post := shared.Post{User: &user, Image: "image.jpg"}
	id, _ := ts.usecase.Create(context.Background(), post)

	post.Image = "updated_image.jpg"
	err := ts.usecase.Update(context.Background(), post, id)
	assert.NoError(t, err)

	updatedPost, err := ts.usecase.GetPost(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, "updated_image.jpg", updatedPost.Image)
}

func TestUpdatePostEmptyImage(t *testing.T) {
	ts := setup()

	user := shared.User{ID: uuid.New(), Username: "testuser"}
	post := shared.Post{User: &user, Image: "image.jpg"}
	id, _ := ts.usecase.Create(context.Background(), post)

	post.Image = ""
	err := ts.usecase.Update(context.Background(), post, id)
	assert.Error(t, err)
}

func TestDeletePost(t *testing.T) {
	ts := setup()

	user := shared.User{ID: uuid.New(), Username: "testuser"}
	post := shared.Post{User: &user, Image: "image.jpg"}
	id, _ := ts.usecase.Create(context.Background(), post)

	err := ts.usecase.Delete(context.Background(), id)
	assert.NoError(t, err)

	posts, err := ts.usecase.GetPosts(context.Background(), 10, time.Now(), uuid.Nil)
	assert.NoError(t, err)
	assert.NotContains(t, posts, post)
}

func TestDeletePostNotFound(t *testing.T) {
	ts := setup()

	id := uuid.New()

	err := ts.usecase.Delete(context.Background(), id)
	assert.Error(t, err)
	assert.Equal(t, "post not found", err.Error())
}

// mockPostRepository is a mock implementation of postRepository for testing
type mockPostRepository struct {
	posts map[uuid.UUID]shared.Post
}

func newMockPostRepository() *mockPostRepository {
	return &mockPostRepository{
		posts: make(map[uuid.UUID]shared.Post),
	}
}

func (m *mockPostRepository) create(ctx context.Context, post shared.Post) (uuid.UUID, error) {
	id := uuid.New()
	m.posts[id] = post

	return id, nil
}

func (m *mockPostRepository) getPosts(ctx context.Context, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]shared.Post, error) {
	var result []shared.Post
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

func (m *mockPostRepository) getPost(ctx context.Context, id uuid.UUID) (shared.Post, error) {
	post, exists := m.posts[id]
	if !exists {
		return shared.Post{}, fmt.Errorf("post not found")
	}

	return post, nil
}

func (m *mockPostRepository) update(ctx context.Context, post shared.Post, id uuid.UUID) error {
	if _, exists := m.posts[id]; !exists {
		return fmt.Errorf("post not found")
	}
	m.posts[id] = post

	return nil
}

func (m *mockPostRepository) delete(ctx context.Context, id uuid.UUID) error {
	if _, exists := m.posts[id]; !exists {
		return fmt.Errorf("post not found")
	}
	delete(m.posts, id)

	return nil
}
