package users

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestSetup struct {
	usecase UserUsecaseI
	repo    *mockUserRepository
}

func setup() *TestSetup {
	repo := newMockUserRepository()
	usecase := &userUsecase{repository: repo}

	return &TestSetup{usecase: usecase, repo: repo}
}

func TestCreateUser(t *testing.T) {
	ts := setup()
	user := User{Username: "testuser", Password: "password123"}

	id, err := ts.usecase.Create(user)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestCreateUser_EmptyFields(t *testing.T) {
	ts := setup()
	user := User{Username: "", Password: ""}

	id, err := ts.usecase.Create(user)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestGetUsers(t *testing.T) {
	ts := setup()
	user1 := User{Username: "testuser", Password: "password123"}
	user2 := User{Username: "testuser2", Password: "password456"}
	id1, _ := ts.usecase.Create(user1)
	id2, _ := ts.usecase.Create(user2)

	users, err := ts.usecase.Get([]uuid.UUID{id1, id2})
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestGetBySearch(t *testing.T) {
	ts := setup()
	user1 := User{Username: "testuser1", Password: "password123"}
	user2 := User{Username: "testuser2", Password: "password456"}
	ts.repo.create(user1)
	ts.repo.create(user2)

	users, err := ts.usecase.GetBySearch("testuser")
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUpdateUser(t *testing.T) {
	ts := setup()
	user := User{Username: "testuser", Password: "password123"}
	id, _ := ts.usecase.Create(user)

	user.Username = "updateduser"
	err := ts.usecase.Update(user, id)
	assert.NoError(t, err)

	updatedUsers, err := ts.usecase.Get([]uuid.UUID{id})
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", updatedUsers[0].Username)
}

func TestUpdateUser_NotFound(t *testing.T) {
	ts := setup()
	user := User{Username: "testuser", Password: "password123"}
	id := uuid.New()

	err := ts.usecase.Update(user, id)
	assert.Error(t, err)
}

func TestDeleteUser(t *testing.T) {
	ts := setup()
	user := User{Username: "testuser", Password: "password123"}
	id, _ := ts.usecase.Create(user)

	err := ts.usecase.Delete(id)
	assert.NoError(t, err)

	users, err := ts.usecase.Get([]uuid.UUID{id})
	assert.NoError(t, err)
	assert.Len(t, users, 0)
}

func TestDeleteUser_NotFound(t *testing.T) {
	ts := setup()
	id := uuid.New()

	err := ts.usecase.Delete(id)
	assert.NoError(t, err)
}

// mockUserRepository is a mock implementation of userRepositoryI for testing
type mockUserRepository struct {
	users map[uuid.UUID]User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[uuid.UUID]User),
	}
}

func (m *mockUserRepository) create(user User) (uuid.UUID, error) {
	id := uuid.New()
	m.users[id] = user

	return id, nil
}

func (m *mockUserRepository) get(idList []uuid.UUID) ([]User, error) {
	var result []User
	for _, id := range idList {
		if user, exists := m.users[id]; exists {
			result = append(result, user)
		}
	}

	return result, nil
}

func (m *mockUserRepository) getBySearch(searchStr string) ([]User, error) {
	var result []User
	for _, user := range m.users {
		if user.Username == searchStr {
			result = append(result, user)
		}
	}

	return result, nil
}

func (m *mockUserRepository) update(user User, id uuid.UUID) error {
	if _, exists := m.users[id]; !exists {
		return fmt.Errorf("user not found")
	}
	m.users[id] = user

	return nil
}

func (m *mockUserRepository) delete(id uuid.UUID) error {
	delete(m.users, id)

	return nil
}
