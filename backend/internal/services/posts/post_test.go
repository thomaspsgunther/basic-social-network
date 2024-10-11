package posts

import (
	"context"
	"fmt"
	"testing"
	"time"
	"y-net/internal/services/shared"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestSetup struct {
	usecase IPostUsecase
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

func TestLike(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()

	err := ts.usecase.Like(context.Background(), userId, postId)
	assert.NoError(t, err)

	likedUsers, err := ts.usecase.GetLikes(context.Background(), postId)
	assert.NoError(t, err)
	assert.Len(t, likedUsers, 1)
	assert.Equal(t, likedUsers[0].ID, userId)
}

func TestUnlike(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()

	err := ts.usecase.Like(context.Background(), userId, postId)
	assert.NoError(t, err)

	err = ts.usecase.Unlike(context.Background(), userId, postId)
	assert.NoError(t, err)

	likedUsers, err := ts.usecase.GetLikes(context.Background(), postId)
	assert.NoError(t, err)
	assert.Len(t, likedUsers, 0)
}

func TestUnlikeUserNotLiked(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()

	err := ts.usecase.Unlike(context.Background(), userId, postId)
	assert.Error(t, err)

	likedUsers, err := ts.usecase.GetLikes(context.Background(), postId)
	assert.NoError(t, err)
	assert.Len(t, likedUsers, 0)
}

func TestUserLikedPost(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()

	err := ts.usecase.Like(context.Background(), userId, postId)
	assert.NoError(t, err)

	liked, err := ts.usecase.UserLikedPost(context.Background(), userId, postId)
	assert.NoError(t, err)
	assert.True(t, liked)

	anotherUserId := uuid.New()
	liked, err = ts.usecase.UserLikedPost(context.Background(), anotherUserId, postId)
	assert.NoError(t, err)
	assert.False(t, liked)

	anotherPostId := uuid.New()
	liked, err = ts.usecase.UserLikedPost(context.Background(), userId, anotherPostId)
	assert.Error(t, err)
	assert.False(t, liked)
}

func TestGetLikes(t *testing.T) {
	ts := setup()

	postId := uuid.New()
	userId1 := uuid.New()
	userId2 := uuid.New()

	err := ts.usecase.Like(context.Background(), userId1, postId)
	assert.NoError(t, err)

	err = ts.usecase.Like(context.Background(), userId2, postId)
	assert.NoError(t, err)

	likedUsers, err := ts.usecase.GetLikes(context.Background(), postId)
	assert.NoError(t, err)
	assert.Len(t, likedUsers, 2)
	assert.Contains(t, likedUsers, shared.User{ID: userId1})
	assert.Contains(t, likedUsers, shared.User{ID: userId2})
}

// mockPostRepository is a mock implementation of iPostRepository for testing
type mockPostRepository struct {
	posts map[uuid.UUID]shared.Post
	likes map[uuid.UUID][]uuid.UUID
}

func newMockPostRepository() *mockPostRepository {
	return &mockPostRepository{
		posts: make(map[uuid.UUID]shared.Post),
		likes: make(map[uuid.UUID][]uuid.UUID),
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

func (m *mockPostRepository) like(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error {
	if m.likes[postId] == nil {
		m.likes[postId] = []uuid.UUID{}
	}
	m.likes[postId] = append(m.likes[postId], userId)

	return nil
}

func (m *mockPostRepository) unlike(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error {
	users := m.likes[postId]
	for i, id := range users {
		if id == userId {
			m.likes[postId] = append(users[:i], users[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("user has not liked this post")
}

func (m *mockPostRepository) userLikedPost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error) {
	users, exists := m.likes[postId]
	if !exists {
		return false, fmt.Errorf("post not found")
	}

	for _, id := range users {
		if id == userId {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockPostRepository) getLikes(ctx context.Context, id uuid.UUID) ([]shared.User, error) {
	var userList []shared.User
	for _, userId := range m.likes[id] {
		userList = append(userList, shared.User{ID: userId})
	}

	return userList, nil
}
