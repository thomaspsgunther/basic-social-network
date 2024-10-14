package comments

import (
	"context"
	"fmt"
	"testing"
	"y-net/internal/services/shared"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestSetup struct {
	usecase ICommentUsecase
	repo    *mockCommentRepository
}

func setup() *TestSetup {
	repo := newMockCommentRepository()
	usecase := &commentUsecaseImpl{repository: repo}

	return &TestSetup{usecase: usecase, repo: repo}
}

func TestCreateComment(t *testing.T) {
	ts := setup()

	user := shared.User{ID: uuid.New(), Username: "testuser"}
	postId := uuid.New()
	comment := Comment{User: &user, PostID: postId, Message: "This is a comment."}

	id, err := ts.usecase.Create(context.Background(), comment)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestCreateCommentEmptyUserAndNilPostID(t *testing.T) {
	ts := setup()

	comment := Comment{User: &shared.User{}, PostID: uuid.Nil, Message: "This is a comment."}

	id, err := ts.usecase.Create(context.Background(), comment)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestCreateCommentEmptyDescription(t *testing.T) {
	ts := setup()

	user := shared.User{ID: uuid.New(), Username: "testuser"}
	postId := uuid.New()
	comment := Comment{User: &user, PostID: postId, Message: ""}

	id, err := ts.usecase.Create(context.Background(), comment)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestUpdateComment(t *testing.T) {
	ts := setup()

	user := shared.User{ID: uuid.New(), Username: "testuser"}
	postId := uuid.New()
	comment := Comment{User: &user, PostID: postId, Message: "This is a comment."}
	id, _ := ts.usecase.Create(context.Background(), comment)

	comment.Message = "Updated comment."
	err := ts.usecase.Update(context.Background(), comment, id)
	assert.NoError(t, err)

	updatedComment, err := ts.usecase.Get(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, "Updated comment.", updatedComment.Message)
}

func TestUpdateCommentNotFound(t *testing.T) {
	ts := setup()

	comment := Comment{Message: "This is a comment."}
	id := uuid.New()

	err := ts.usecase.Update(context.Background(), comment, id)
	assert.Error(t, err)
	assert.Equal(t, "comment not found", err.Error())
}

func TestLikeComment(t *testing.T) {
	ts := setup()

	user := shared.User{ID: uuid.New(), Username: "testuser"}
	postId := uuid.New()
	comment := Comment{User: &user, PostID: postId, Message: "This is a comment."}
	id, _ := ts.usecase.Create(context.Background(), comment)

	err := ts.usecase.Like(context.Background(), id)
	assert.NoError(t, err)

	likedComment, err := ts.usecase.Get(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, 1, likedComment.LikeCount)
}

func TestLikeCommentNotFound(t *testing.T) {
	ts := setup()

	id := uuid.New()

	err := ts.usecase.Like(context.Background(), id)
	assert.Error(t, err)
	assert.Equal(t, "comment not found", err.Error())
}

func TestUnlikeComment(t *testing.T) {
	ts := setup()

	user := shared.User{ID: uuid.New(), Username: "testuser"}
	postId := uuid.New()
	comment := Comment{User: &user, PostID: postId, Message: "This is a comment."}
	id, _ := ts.usecase.Create(context.Background(), comment)

	_ = ts.usecase.Like(context.Background(), id)
	err := ts.usecase.Unlike(context.Background(), id)
	assert.NoError(t, err)

	unlikedComment, err := ts.usecase.Get(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, 0, unlikedComment.LikeCount)
}

func TestUnlikeCommentNotFound(t *testing.T) {
	ts := setup()

	id := uuid.New()

	err := ts.usecase.Unlike(context.Background(), id)
	assert.Error(t, err)
	assert.Equal(t, "comment not found", err.Error())
}

func TestDeleteComment(t *testing.T) {
	ts := setup()

	user := shared.User{ID: uuid.New(), Username: "testuser"}
	postId := uuid.New()
	comment := Comment{User: &user, PostID: postId, Message: "This is a comment."}
	id, _ := ts.usecase.Create(context.Background(), comment)

	err := ts.usecase.Delete(context.Background(), id)
	assert.NoError(t, err)

	comments, err := ts.usecase.GetFromPost(context.Background(), comment.PostID)
	assert.NoError(t, err)
	assert.NotContains(t, comments, comment)
}

func TestDeleteCommentNotFound(t *testing.T) {
	ts := setup()

	id := uuid.New()

	err := ts.usecase.Delete(context.Background(), id)
	assert.Error(t, err)
	assert.Equal(t, "comment not found", err.Error())
}

func TestGetFromPost(t *testing.T) {
	ts := setup()

	postID := uuid.New()
	comment1 := Comment{Message: "First comment", PostID: postID}
	comment2 := Comment{Message: "Second comment", PostID: postID}
	_, _ = ts.usecase.Create(context.Background(), comment1)
	_, _ = ts.usecase.Create(context.Background(), comment2)

	comments, err := ts.usecase.GetFromPost(context.Background(), postID)
	assert.NoError(t, err)
	assert.Len(t, comments, 2)
}

func TestGetComment(t *testing.T) {
	ts := setup()

	user := shared.User{ID: uuid.New(), Username: "testuser"}
	postId := uuid.New()
	comment := Comment{User: &user, PostID: postId, Message: "This is a comment."}
	id, _ := ts.usecase.Create(context.Background(), comment)

	retrievedComment, err := ts.usecase.Get(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, comment.Message, retrievedComment.Message)
}

func TestGetCommentNotFound(t *testing.T) {
	ts := setup()

	id := uuid.New()

	_, err := ts.usecase.Get(context.Background(), id)
	assert.Error(t, err)
}

// mockCommentRepository is a mock implementation of iCommentRepository for testing purposes
type mockCommentRepository struct {
	comments map[uuid.UUID]Comment
}

func newMockCommentRepository() *mockCommentRepository {
	return &mockCommentRepository{
		comments: make(map[uuid.UUID]Comment),
	}
}

func (m *mockCommentRepository) create(ctx context.Context, comment Comment) (uuid.UUID, error) {
	id := uuid.New()
	comment.ID = id
	m.comments[id] = comment

	return id, nil
}

func (m *mockCommentRepository) update(ctx context.Context, comment Comment, id uuid.UUID) error {
	if _, exists := m.comments[id]; !exists {
		return fmt.Errorf("comment not found")
	}
	comment.ID = id
	m.comments[id] = comment

	return nil
}

func (m *mockCommentRepository) like(ctx context.Context, id uuid.UUID) error {
	if comment, exists := m.comments[id]; exists {
		comment.LikeCount++
		m.comments[id] = comment

		return nil
	}

	return fmt.Errorf("comment not found")
}

func (m *mockCommentRepository) unlike(ctx context.Context, id uuid.UUID) error {
	if comment, exists := m.comments[id]; exists {
		if comment.LikeCount > 0 {
			comment.LikeCount--
		}
		m.comments[id] = comment

		return nil
	}

	return fmt.Errorf("comment not found")
}

func (m *mockCommentRepository) delete(ctx context.Context, id uuid.UUID) error {
	if _, exists := m.comments[id]; !exists {
		return fmt.Errorf("comment not found")
	}
	delete(m.comments, id)

	return nil
}

func (m *mockCommentRepository) get(ctx context.Context, id uuid.UUID) (Comment, error) {
	if comment, exists := m.comments[id]; exists {
		return comment, nil
	}

	return Comment{}, fmt.Errorf("comment not found")
}

func (m *mockCommentRepository) getFromPost(ctx context.Context, postId uuid.UUID) ([]Comment, error) {
	var result []Comment
	for _, comment := range m.comments {
		if comment.PostID == postId {
			result = append(result, comment)
		}
	}

	return result, nil
}
