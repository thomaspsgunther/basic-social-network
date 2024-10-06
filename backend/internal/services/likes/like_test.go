package likes

import (
	"fmt"
	"testing"
	"y_net/internal/services/users"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestSetup struct {
	usecase LikeUsecaseI
	repo    *mockLikeRepository
}

func setup() *TestSetup {
	repo := newMockLikeRepository()
	usecase := &likeUsecase{repository: repo}

	return &TestSetup{usecase: usecase, repo: repo}
}

func TestLikePost(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()

	err := ts.usecase.LikePost(userId, postId)
	assert.NoError(t, err)

	likedUsers, err := ts.usecase.GetFromPost(postId)
	assert.NoError(t, err)
	assert.Len(t, likedUsers, 1)
	assert.Equal(t, likedUsers[0].ID, userId)
}

func TestUnlikePost(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()

	err := ts.usecase.LikePost(userId, postId)
	assert.NoError(t, err)

	err = ts.usecase.UnlikePost(userId, postId)
	assert.NoError(t, err)

	likedUsers, err := ts.usecase.GetFromPost(postId)
	assert.NoError(t, err)
	assert.Len(t, likedUsers, 0)
}

func TestUnlikePostUserNotLiked(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()

	err := ts.usecase.UnlikePost(userId, postId)
	assert.Error(t, err)

	likedUsers, err := ts.usecase.GetFromPost(postId)
	assert.NoError(t, err)
	assert.Len(t, likedUsers, 0)
}

func TestGetFromPost(t *testing.T) {
	ts := setup()

	postId := uuid.New()
	userId1 := uuid.New()
	userId2 := uuid.New()

	err := ts.usecase.LikePost(userId1, postId)
	assert.NoError(t, err)

	err = ts.usecase.LikePost(userId2, postId)
	assert.NoError(t, err)

	likedUsers, err := ts.usecase.GetFromPost(postId)
	assert.NoError(t, err)
	assert.Len(t, likedUsers, 2)
	assert.Contains(t, likedUsers, users.User{ID: userId1})
	assert.Contains(t, likedUsers, users.User{ID: userId2})
}

// mockLikeRepository is a mock implementation of likeRepositoryI for testing purposes
type mockLikeRepository struct {
	likes map[uuid.UUID][]uuid.UUID
}

func newMockLikeRepository() *mockLikeRepository {
	return &mockLikeRepository{
		likes: make(map[uuid.UUID][]uuid.UUID),
	}
}

func (m *mockLikeRepository) likePost(userId uuid.UUID, postId uuid.UUID) error {
	if m.likes[postId] == nil {
		m.likes[postId] = []uuid.UUID{}
	}
	m.likes[postId] = append(m.likes[postId], userId)

	return nil
}

func (m *mockLikeRepository) unlikePost(userId uuid.UUID, postId uuid.UUID) error {
	users := m.likes[postId]
	for i, id := range users {
		if id == userId {
			m.likes[postId] = append(users[:i], users[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("user has not liked this post")
}

func (m *mockLikeRepository) getFromPost(postId uuid.UUID) ([]users.User, error) {
	var userList []users.User
	for _, userId := range m.likes[postId] {
		userList = append(userList, users.User{ID: userId})
	}

	return userList, nil
}
