package users

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
	usecase IUserUsecase
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

func TestGetUser(t *testing.T) {
	ts := setup()

	user1 := shared.User{Username: "testuser", Password: "password123"}
	id1, _ := ts.usecase.Create(context.Background(), user1)

	user, err := ts.usecase.Get(context.Background(), id1)
	assert.NoError(t, err)
	assert.Equal(t, user1.Username, user.Username)
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

	updatedUser, err := ts.usecase.Get(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", updatedUser.Username)
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

	_, err = ts.usecase.Get(context.Background(), id)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestDeleteUserNotFound(t *testing.T) {
	ts := setup()

	id := uuid.New()

	err := ts.usecase.Delete(context.Background(), id)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestFollow(t *testing.T) {
	ts := setup()

	followerId := uuid.New()
	followedId := uuid.New()

	err := ts.usecase.Follow(context.Background(), followerId, followedId)
	assert.NoError(t, err)

	followedUsers, err := ts.usecase.GetFollowed(context.Background(), followerId)
	assert.NoError(t, err)
	assert.Len(t, followedUsers, 1)
	assert.Equal(t, followedUsers[0].ID, followedId)
}

func TestUnfollow(t *testing.T) {
	ts := setup()

	followerId := uuid.New()
	followedId := uuid.New()

	err := ts.usecase.Follow(context.Background(), followerId, followedId)
	assert.NoError(t, err)

	err = ts.usecase.Unfollow(context.Background(), followerId, followedId)
	assert.NoError(t, err)

	followedUsers, err := ts.usecase.GetFollowed(context.Background(), followerId)
	assert.NoError(t, err)
	assert.Len(t, followedUsers, 0)
}

func TestUnfollow_NonExistent(t *testing.T) {
	ts := setup()

	followerId := uuid.New()
	nonExistentUserId := uuid.New()

	err := ts.usecase.Unfollow(context.Background(), followerId, nonExistentUserId)
	assert.Error(t, err)
	assert.Equal(t, "user "+nonExistentUserId.String()+" is not being followed", err.Error())
}

func TestIsFollowing(t *testing.T) {
	ts := setup()

	followerId := uuid.New()
	followedId := uuid.New()

	err := ts.usecase.Follow(context.Background(), followerId, followedId)
	assert.NoError(t, err)

	isFollowing, err := ts.usecase.UserFollowsUser(context.Background(), followerId, followedId)
	assert.NoError(t, err)
	assert.True(t, isFollowing)
}

func TestGetFollowers(t *testing.T) {
	ts := setup()

	followerId := uuid.New()
	followedId := uuid.New()

	err := ts.usecase.Follow(context.Background(), followerId, followedId)
	assert.NoError(t, err)

	followersList, err := ts.usecase.GetFollowers(context.Background(), followedId)
	assert.NoError(t, err)
	assert.Len(t, followersList, 1)
	assert.Equal(t, followersList[0].ID, followerId)
}

func TestGetFollowed(t *testing.T) {
	ts := setup()

	followerId := uuid.New()
	followedId1 := uuid.New()
	followedId2 := uuid.New()

	err := ts.usecase.Follow(context.Background(), followerId, followedId1)
	assert.NoError(t, err)
	err = ts.usecase.Follow(context.Background(), followerId, followedId2)
	assert.NoError(t, err)

	followedUsers, err := ts.usecase.GetFollowed(context.Background(), followerId)
	assert.NoError(t, err)
	assert.Len(t, followedUsers, 2)
	assert.Contains(t, followedUsers, shared.User{ID: followedId1})
	assert.Contains(t, followedUsers, shared.User{ID: followedId2})
}

// mockUserRepository is a mock implementation of iUserRepository for testing
type mockUserRepository struct {
	users        map[uuid.UUID]shared.User
	followersMap map[uuid.UUID][]uuid.UUID
	posts        map[uuid.UUID]shared.Post
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:        make(map[uuid.UUID]shared.User),
		followersMap: make(map[uuid.UUID][]uuid.UUID),
		posts:        make(map[uuid.UUID]shared.Post),
	}
}

func (m *mockUserRepository) create(ctx context.Context, user shared.User) (uuid.UUID, error) {
	id := uuid.New()
	m.users[id] = user

	return id, nil
}

func (m *mockUserRepository) get(ctx context.Context, id uuid.UUID) (shared.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return shared.User{}, fmt.Errorf("user not found")
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

func (m *mockUserRepository) follow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error {
	m.followersMap[followerId] = append(m.followersMap[followerId], followedId)

	return nil
}

func (m *mockUserRepository) getFollowers(ctx context.Context, userId uuid.UUID) ([]shared.User, error) {
	var userList []shared.User
	for followerId := range m.followersMap {
		for _, id := range m.followersMap[followerId] {
			if id == userId {
				userList = append(userList, shared.User{ID: followerId})
			}
		}
	}

	return userList, nil
}

func (m *mockUserRepository) getFollowed(ctx context.Context, id uuid.UUID) ([]shared.User, error) {
	var userList []shared.User
	for followerId, followed := range m.followersMap {
		for _, id := range followed {
			if followerId == id {
				userList = append(userList, shared.User{ID: id})
			}
		}
	}

	return userList, nil
}

func (m *mockUserRepository) unfollow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error {
	followed := m.followersMap[followerId]
	for i, id := range followed {
		if id == followedId {
			m.followersMap[followerId] = append(followed[:i], followed[i+1:]...)

			return nil
		}
	}

	return fmt.Errorf("user %v is not being followed", followedId)
}

func (m *mockUserRepository) userFollowsUser(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) (bool, error) {
	followed := m.followersMap[followerId]
	if followed == nil {
		return false, nil
	}

	for _, id := range followed {
		if id == followedId {
			return true, nil
		}
	}

	return false, nil
}
