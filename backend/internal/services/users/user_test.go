package users

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
	usecase UserUsecase
	repo    *mockUserRepository
}

func setup() *TestSetup {
	repo := newMockUserRepository()
	usecase := &userUsecaseImpl{repository: repo}

	return &TestSetup{usecase: usecase, repo: repo}
}

func TestCreateUser(t *testing.T) {
	ts := setup()

	user := shared.User{Username: "testuser", Password: "password123"}

	id, err := ts.usecase.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestCreateUserEmptyFields(t *testing.T) {
	ts := setup()

	user := shared.User{Username: "", Password: ""}

	id, err := ts.usecase.Create(context.Background(), user)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestGetUsers(t *testing.T) {
	ts := setup()

	user1 := shared.User{Username: "testuser", Password: "password123"}
	user2 := shared.User{Username: "testuser2", Password: "password456"}
	id1, _ := ts.usecase.Create(context.Background(), user1)
	id2, _ := ts.usecase.Create(context.Background(), user2)

	users, err := ts.usecase.Get(context.Background(), []uuid.UUID{id1, id2})
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestGetBySearch(t *testing.T) {
	ts := setup()

	user1 := shared.User{Username: "testuser1", Password: "password123"}
	user2 := shared.User{Username: "testuser2", Password: "password456"}
	ts.usecase.Create(context.Background(), user1)
	ts.usecase.Create(context.Background(), user2)

	users, err := ts.usecase.GetBySearch(context.Background(), "testuser")
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUpdateUser(t *testing.T) {
	ts := setup()

	user := shared.User{Username: "testuser", Password: "password123"}
	id, _ := ts.usecase.Create(context.Background(), user)

	user.Username = "updateduser"
	err := ts.usecase.Update(context.Background(), user, id)
	assert.NoError(t, err)

	updatedUsers, err := ts.usecase.Get(context.Background(), []uuid.UUID{id})
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", updatedUsers[0].Username)
}

func TestUpdateUserNotFound(t *testing.T) {
	ts := setup()

	user := shared.User{Username: "testuser", Password: "password123"}
	id := uuid.New()

	err := ts.usecase.Update(context.Background(), user, id)
	assert.Error(t, err)
}

func TestDeleteUser(t *testing.T) {
	ts := setup()

	user := shared.User{Username: "testuser", Password: "password123"}
	id, _ := ts.usecase.Create(context.Background(), user)

	err := ts.usecase.Delete(context.Background(), id)
	assert.NoError(t, err)

	users, err := ts.usecase.Get(context.Background(), []uuid.UUID{id})
	assert.NoError(t, err)
	assert.Len(t, users, 0)
}

func TestDeleteUserNotFound(t *testing.T) {
	ts := setup()

	id := uuid.New()

	err := ts.usecase.Delete(context.Background(), id)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

// mockUserRepository is a mock implementation of userRepository for testing
type mockUserRepository struct {
	users map[uuid.UUID]shared.User
	posts map[uuid.UUID]shared.Post
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[uuid.UUID]shared.User),
	}
}

func (m *mockUserRepository) create(ctx context.Context, user shared.User) (uuid.UUID, error) {
	id := uuid.New()
	m.users[id] = user

	return id, nil
}

func (m *mockUserRepository) get(ctx context.Context, idList []uuid.UUID) ([]shared.User, error) {
	var result []shared.User
	for _, id := range idList {
		if user, exists := m.users[id]; exists {
			result = append(result, user)
		}
	}

	return result, nil
}

func (m *mockUserRepository) getBySearch(ctx context.Context, searchStr string) ([]shared.User, error) {
	var result []shared.User
	for _, user := range m.users {
		if user.Username == searchStr {
			result = append(result, user)
		}
	}

	return result, nil
}

func (m *mockUserRepository) getPostsFromUser(ctx context.Context, userId uuid.UUID, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]shared.Post, error) {
	var result []shared.Post
	for _, post := range m.posts {
		if len(result) >= limit {
			break
		}
		if post.User.ID == userId && post.CreatedAt.After(lastCreatedAt) && post.ID != lastId {
			result = append(result, post)
		}
	}

	return result, nil
}

func (m *mockUserRepository) update(ctx context.Context, user shared.User, id uuid.UUID) error {
	if _, exists := m.users[id]; !exists {
		return fmt.Errorf("user not found")
	}
	m.users[id] = user

	return nil
}

func (m *mockUserRepository) delete(ctx context.Context, id uuid.UUID) error {
	if _, exists := m.users[id]; !exists {
		return fmt.Errorf("user not found")
	}
	delete(m.users, id)

	return nil
}
