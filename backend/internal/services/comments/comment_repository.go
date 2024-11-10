package comments

import (
	"context"
	"fmt"
	database "y-net/internal/database/postgres"
	"y-net/internal/services/shared"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type iCommentRepository interface {
	create(ctx context.Context, comment Comment) (Comment, error)
	getFromPost(ctx context.Context, postId uuid.UUID) ([]Comment, error)
	get(ctx context.Context, id uuid.UUID) (Comment, error)
	update(ctx context.Context, comment Comment, id uuid.UUID) error
	delete(ctx context.Context, id uuid.UUID) error
}

type commentRepositoryImpl struct{}

func (r *commentRepositoryImpl) create(ctx context.Context, comment Comment) (Comment, error) {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return Comment{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	var newComment Comment
	err = tx.QueryRow(
		ctx,
		"INSERT INTO comments (user_id, post_id, message) VALUES ($1, $2, $3) RETURNING id, created_at",
		comment.User.ID, comment.PostID, comment.Message,
	).Scan(&newComment.ID, &newComment.CreatedAt)
	if err != nil {
		return Comment{}, fmt.Errorf("failed to insert comment: %w", err)
	}

	return newComment, nil
}

func (r *commentRepositoryImpl) getFromPost(ctx context.Context, postId uuid.UUID) ([]Comment, error) {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	query := `
		SELECT c.id, c.user_id, u.username, u.avatar, c.post_id, c.message, c.created_at
		FROM comments c
		INNER JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`

	rows, err := tx.Query(ctx, query, postId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []Comment{}, nil
		}

		return nil, fmt.Errorf("failed to select comments: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		comment.User = &shared.User{}
		err := rows.Scan(&comment.ID, &comment.User.ID, &comment.User.Username, &comment.User.Avatar, &comment.PostID, &comment.Message, &comment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return comments, nil
}

func (r *commentRepositoryImpl) get(ctx context.Context, id uuid.UUID) (Comment, error) {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return Comment{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	var comment Comment
	comment.User = &shared.User{}
	err = tx.QueryRow(
		ctx,
		"SELECT id, user_id, post_id FROM comments WHERE id = $1",
		id).Scan(&comment.ID, &comment.User.ID, &comment.PostID)
	if err != nil {
		return Comment{}, fmt.Errorf("failed to scan comment: %w", err)
	}

	return comment, nil
}

func (r *commentRepositoryImpl) update(ctx context.Context, comment Comment, id uuid.UUID) error {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	_, err = tx.Exec(
		ctx,
		"UPDATE comments SET message = $1 WHERE id = $1",
		comment.Message, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

func (r *commentRepositoryImpl) delete(ctx context.Context, id uuid.UUID) error {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	_, err = tx.Exec(ctx, "DELETE FROM comments WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
