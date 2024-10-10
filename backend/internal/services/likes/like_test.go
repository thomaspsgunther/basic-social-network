package likes

import (
	"context"
	"fmt"
	"testing"
	"y_net/internal/services/shared"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestSetup struct {
	usecase LikeUsecase
	repo    *mockLikeRepository
}

func setup() *TestSetup {
	repo := newMockLikeRepository()
	usecase := &likeUsecaseImpl{repository: repo}

	return &TestSetup{usecase: usecase, repo: repo}
}

func TestLikePost(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()

	err := ts.usecase.LikePost(context.Background(), userId, postId)
	assert.NoError(t, err)

	likedUsers, err := ts.usecase.GetFromPost(context.Background(), postId)
	assert.NoError(t, err)
	assert.Len(t, likedUsers, 1)
	assert.Equal(t, likedUsers[0].ID, userId)
}

func TestUnlikePost(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()

	err := ts.usecase.LikePost(context.Background(), userId, postId)
	assert.NoError(t, err)

	err = ts.usecase.UnlikePost(context.Background(), userId, postId)
	assert.NoError(t, err)

	likedUsers, err := ts.usecase.GetFromPost(context.Background(), postId)
	assert.NoError(t, err)
	assert.Len(t, likedUsers, 0)
}

func TestUnlikePostUserNotLiked(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()

	err := ts.usecase.UnlikePost(context.Background(), userId, postId)
	assert.Error(t, err)

	likedUsers, err := ts.usecase.GetFromPost(context.Background(), postId)
	assert.NoError(t, err)
	assert.Len(t, likedUsers, 0)
}

func TestUserLikedPost(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()

	err := ts.usecase.LikePost(context.Background(), userId, postId)
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

func TestGetFromPost(t *testing.T) {
	ts := setup()

	postId := uuid.New()
	userId1 := uuid.New()
	userId2 := uuid.New()

	err := ts.usecase.LikePost(context.Background(), userId1, postId)
	assert.NoError(t, err)

	err = ts.usecase.LikePost(context.Background(), userId2, postId)
	assert.NoError(t, err)

	likedUsers, err := ts.usecase.GetFromPost(context.Background(), postId)
	assert.NoError(t, err)
	assert.Len(t, likedUsers, 2)
	assert.Contains(t, likedUsers, shared.User{ID: userId1})
	assert.Contains(t, likedUsers, shared.User{ID: userId2})
}

// mockLikeRepository is a mock implementation of likeRepository for testing purposes
type mockLikeRepository struct {
	likes map[uuid.UUID][]uuid.UUID
}

func newMockLikeRepository() *mockLikeRepository {
	return &mockLikeRepository{
		likes: make(map[uuid.UUID][]uuid.UUID),
	}
}

func (m *mockLikeRepository) likePost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error {
	if m.likes[postId] == nil {
		m.likes[postId] = []uuid.UUID{}
	}
	m.likes[postId] = append(m.likes[postId], userId)

	return nil
}

func (m *mockLikeRepository) unlikePost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error {
	users := m.likes[postId]
	for i, id := range users {
		if id == userId {
			m.likes[postId] = append(users[:i], users[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("user has not liked this post")
}

func (m *mockLikeRepository) userLikedPost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error) {
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

func (m *mockLikeRepository) getFromPost(ctx context.Context, postId uuid.UUID) ([]shared.User, error) {
	var userList []shared.User
	for _, userId := range m.likes[postId] {
		userList = append(userList, shared.User{ID: userId})
	}

	return userList, nil
}
