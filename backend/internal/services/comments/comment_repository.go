package comments

import (
	"context"
	"fmt"
	database "y_net/internal/database/postgres"

	"github.com/google/uuid"
)

type iCommentRepository interface {
	create(ctx context.Context, comment Comment) (uuid.UUID, error)
	getFromPost(ctx context.Context, postId uuid.UUID) ([]Comment, error)
	get(ctx context.Context, id uuid.UUID) (Comment, error)
	update(ctx context.Context, comment Comment, id uuid.UUID) error
	like(ctx context.Context, id uuid.UUID) error
	unlike(ctx context.Context, id uuid.UUID) error
	delete(ctx context.Context, id uuid.UUID) error
}

type commentRepositoryImpl struct{}

func (r *commentRepositoryImpl) create(ctx context.Context, comment Comment) (uuid.UUID, error) {
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	var id uuid.UUID
	err = tx.QueryRow(
		ctx,
		"INSERT INTO comments (user_id, post_id, description) VALUES ($1, $2, $3) RETURNING id",
		comment.User.ID, comment.PostID, comment.Description,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert comment: %w", err)
	}

	return id, nil
}

func (r *commentRepositoryImpl) getFromPost(ctx context.Context, postId uuid.UUID) ([]Comment, error) {
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	query := `
		SELECT c.id, c.user_id, u.username, u.avatar, c.image, c.description, c.like_count, c.created_at
		FROM comments c
		INNER JOIN users u ON c.user_id = u.id
		WHERE post_id = $1
		ORDER BY like_count DESC
	`

	rows, err := tx.Query(ctx, query, postId)
	if err != nil {
		return nil, fmt.Errorf("failed to select comments: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.User.ID, &comment.User.Username, &comment.User.Avatar, &comment.Description, &comment.LikeCount, &comment.CreatedAt)
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
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return Comment{}, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return Comment{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	var comment Comment
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
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	_, err = tx.Exec(
		ctx,
		"UPDATE comments SET description = $1 WHERE id = $1",
		comment.Description, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

func (r *commentRepositoryImpl) like(ctx context.Context, id uuid.UUID) error {
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	_, err = tx.Exec(
		ctx,
		"UPDATE comments SET like_count = like_count + 1 WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

func (r *commentRepositoryImpl) unlike(ctx context.Context, id uuid.UUID) error {
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	_, err = tx.Exec(
		ctx,
		"UPDATE comments SET like_count = like_count - 1 WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

func (r *commentRepositoryImpl) delete(ctx context.Context, id uuid.UUID) error {
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	_, err = tx.Exec(ctx, "DELETE FROM comments WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}
