package followers

import (
	"fmt"
	"testing"
	"y_net/internal/services/users"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestSetup struct {
	usecase FollowerUsecaseI
	repo    *mockFollowerRepository
}

func setup() *TestSetup {
	repo := newMockFollowerRepository()
	usecase := &followerUsecase{repository: repo}

	return &TestSetup{usecase: usecase, repo: repo}
}

func TestFollow(t *testing.T) {
	ts := setup()

	followerId := uuid.New()
	followedId := uuid.New()

	err := ts.usecase.Follow(followerId, followedId)
	assert.NoError(t, err)

	followedUsers, err := ts.repo.getFollowed(followerId)
	assert.NoError(t, err)
	assert.Len(t, followedUsers, 1)
	assert.Equal(t, followedUsers[0].ID, followedId)
}

func TestUnfollow(t *testing.T) {
	ts := setup()

	followerId := uuid.New()
	followedId := uuid.New()

	err := ts.usecase.Follow(followerId, followedId)
	assert.NoError(t, err)

	err = ts.usecase.Unfollow(followerId, followedId)
	assert.NoError(t, err)

	followedUsers, err := ts.repo.getFollowed(followerId)
	assert.NoError(t, err)
	assert.Len(t, followedUsers, 0)
}

func TestUnfollow_NonExistent(t *testing.T) {
	ts := setup()

	followerId := uuid.New()
	nonExistentUserId := uuid.New()

	err := ts.usecase.Unfollow(followerId, nonExistentUserId)
	assert.Error(t, err)
	assert.Equal(t, "user "+nonExistentUserId.String()+" is not being followed", err.Error())
}

func TestIsFollowing(t *testing.T) {
	ts := setup()

	followerId := uuid.New()
	followedId := uuid.New()

	err := ts.usecase.Follow(followerId, followedId)
	assert.NoError(t, err)

	isFollowing, err := ts.usecase.UserFollowsUser(followerId, followedId)
	assert.NoError(t, err)
	assert.True(t, isFollowing)
}

func TestGetFollowers(t *testing.T) {
	ts := setup()

	followerId := uuid.New()
	followedId := uuid.New()

	err := ts.usecase.Follow(followerId, followedId)
	assert.NoError(t, err)

	followersList, err := ts.usecase.GetFollowers(followedId)
	assert.NoError(t, err)
	assert.Len(t, followersList, 1)
	assert.Equal(t, followersList[0].ID, followerId)
}

func TestGetFollowed(t *testing.T) {
	ts := setup()

	followerId := uuid.New()
	followedId1 := uuid.New()
	followedId2 := uuid.New()

	err := ts.usecase.Follow(followerId, followedId1)
	assert.NoError(t, err)
	err = ts.usecase.Follow(followerId, followedId2)
	assert.NoError(t, err)

	followedUsers, err := ts.usecase.GetFollowed(followerId)
	assert.NoError(t, err)
	assert.Len(t, followedUsers, 2)
	assert.Contains(t, followedUsers, users.User{ID: followedId1})
	assert.Contains(t, followedUsers, users.User{ID: followedId2})
}

// mockFollowerRepository is a mock implementation of followerRepositoryI for testing purposes
type mockFollowerRepository struct {
	followersMap map[uuid.UUID][]uuid.UUID
}

func newMockFollowerRepository() *mockFollowerRepository {
	return &mockFollowerRepository{
		followersMap: make(map[uuid.UUID][]uuid.UUID),
	}
}

func (m *mockFollowerRepository) follow(followerId uuid.UUID, followedId uuid.UUID) error {
	m.followersMap[followerId] = append(m.followersMap[followerId], followedId)

	return nil
}

func (m *mockFollowerRepository) unfollow(followerId uuid.UUID, followedId uuid.UUID) error {
	followed := m.followersMap[followerId]
	for i, id := range followed {
		if id == followedId {
			m.followersMap[followerId] = append(followed[:i], followed[i+1:]...)

			return nil
		}
	}

	return fmt.Errorf("user %v is not being followed", followedId)
}

func (m *mockFollowerRepository) userFollowsUser(followerId uuid.UUID, followedId uuid.UUID) (bool, error) {
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

func (m *mockFollowerRepository) getFollowers(userId uuid.UUID) ([]users.User, error) {
	var userList []users.User
	for followerId := range m.followersMap {
		for _, id := range m.followersMap[followerId] {
			if id == userId {
				userList = append(userList, users.User{ID: followerId})
			}
		}
	}

	return userList, nil
}

func (m *mockFollowerRepository) getFollowed(userId uuid.UUID) ([]users.User, error) {
	var userList []users.User
	for followerId, followed := range m.followersMap {
		for _, id := range followed {
			if followerId == userId {
				userList = append(userList, users.User{ID: id})
			}
		}
	}

	return userList, nil
}
