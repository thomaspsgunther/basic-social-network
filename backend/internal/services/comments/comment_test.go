package comments

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestSetup struct {
	usecase CommentUsecaseI
	repo    *mockCommentRepository
}

func setup() *TestSetup {
	repo := newMockCommentRepository()
	usecase := &commentUsecase{repository: repo}

	return &TestSetup{usecase: usecase, repo: repo}
}

func TestCreateComment(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()
	comment := Comment{UserID: userId, PostID: postId, Description: "This is a comment."}

	id, err := ts.usecase.Create(comment)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestCreatePostWithNilUserIDPostID(t *testing.T) {
	ts := setup()

	comment := Comment{UserID: uuid.Nil, PostID: uuid.Nil, Description: "This is a comment."}

	id, err := ts.usecase.Create(comment)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestCreateCommentEmptyDescription(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()
	comment := Comment{UserID: userId, PostID: postId, Description: ""}

	id, err := ts.usecase.Create(comment)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestUpdateComment(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()
	comment := Comment{UserID: userId, PostID: postId, Description: "This is a comment."}
	id, _ := ts.repo.create(comment)

	comment.Description = "Updated comment."
	err := ts.usecase.Update(comment, id)
	assert.NoError(t, err)

	updatedComment, err := ts.usecase.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, "Updated comment.", updatedComment.Description)
}

func TestLikeComment(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()
	comment := Comment{UserID: userId, PostID: postId, Description: "This is a comment."}
	id, _ := ts.repo.create(comment)

	err := ts.usecase.Like(id)
	assert.NoError(t, err)

	likedComment, err := ts.usecase.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, 1, likedComment.LikeCount)
}

func TestUnlikeComment(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()
	comment := Comment{UserID: userId, PostID: postId, Description: "This is a comment."}
	id, _ := ts.repo.create(comment)

	_ = ts.usecase.Like(id)
	err := ts.usecase.Unlike(id)
	assert.NoError(t, err)

	unlikedComment, err := ts.usecase.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, 0, unlikedComment.LikeCount)
}

func TestDeleteComment(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()
	comment := Comment{UserID: userId, PostID: postId, Description: "This is a comment."}
	id, _ := ts.repo.create(comment)

	err := ts.usecase.Delete(id)
	assert.NoError(t, err)

	comments, err := ts.usecase.GetFromPost(comment.PostID)
	assert.NoError(t, err)
	assert.NotContains(t, comments, comment)
}

func TestGetFromPost(t *testing.T) {
	ts := setup()

	postID := uuid.New()
	comment1 := Comment{Description: "First comment", PostID: postID}
	comment2 := Comment{Description: "Second comment", PostID: postID}
	_, _ = ts.repo.create(comment1)
	_, _ = ts.repo.create(comment2)

	comments, err := ts.usecase.GetFromPost(postID)
	assert.NoError(t, err)
	assert.Len(t, comments, 2)
}

func TestGetComment(t *testing.T) {
	ts := setup()

	userId := uuid.New()
	postId := uuid.New()
	comment := Comment{UserID: userId, PostID: postId, Description: "This is a comment."}
	id, _ := ts.repo.create(comment)

	retrievedComment, err := ts.usecase.Get(id)
	assert.NoError(t, err)
	assert.Equal(t, comment.Description, retrievedComment.Description)
}

func TestGetCommentNotFound(t *testing.T) {
	ts := setup()

	id := uuid.New()

	_, err := ts.usecase.Get(id)
	assert.Error(t, err)
}

// mockCommentRepository is a mock implementation of commentRepositoryI for testing purposes
type mockCommentRepository struct {
	comments map[uuid.UUID]Comment
}

func newMockCommentRepository() *mockCommentRepository {
	return &mockCommentRepository{
		comments: make(map[uuid.UUID]Comment),
	}
}

func (m *mockCommentRepository) create(comment Comment) (uuid.UUID, error) {
	id := uuid.New()
	comment.ID = id
	m.comments[id] = comment

	return id, nil
}

func (m *mockCommentRepository) update(comment Comment, id uuid.UUID) error {
	if _, exists := m.comments[id]; !exists {
		return fmt.Errorf("comment not found")
	}
	comment.ID = id
	m.comments[id] = comment

	return nil
}

func (m *mockCommentRepository) like(id uuid.UUID) error {
	if comment, exists := m.comments[id]; exists {
		comment.LikeCount++
		m.comments[id] = comment

		return nil
	}

	return fmt.Errorf("comment not found")
}

func (m *mockCommentRepository) unlike(id uuid.UUID) error {
	if comment, exists := m.comments[id]; exists {
		if comment.LikeCount > 0 {
			comment.LikeCount--
		}
		m.comments[id] = comment

		return nil
	}

	return fmt.Errorf("comment not found")
}

func (m *mockCommentRepository) delete(id uuid.UUID) error {
	if _, exists := m.comments[id]; !exists {
		return fmt.Errorf("comment not found")
	}
	delete(m.comments, id)

	return nil
}

func (m *mockCommentRepository) get(id uuid.UUID) (Comment, error) {
	if comment, exists := m.comments[id]; exists {
		return comment, nil
	}

	return Comment{}, fmt.Errorf("comment not found")
}

func (m *mockCommentRepository) getFromPost(postId uuid.UUID) ([]Comment, error) {
	var result []Comment
	for _, comment := range m.comments {
		if comment.PostID == postId {
			result = append(result, comment)
		}
	}

	return result, nil
}
