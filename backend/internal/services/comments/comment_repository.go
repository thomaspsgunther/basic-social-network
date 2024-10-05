package comments

import (
	"context"
	"fmt"
	database "y_net/internal/database/postgres"

	"github.com/google/uuid"
)

type commentRepositoryI interface {
	create(comment Comment) (uuid.UUID, error)
	update(comment Comment, id uuid.UUID) error
	like(id uuid.UUID) error
	unlike(id uuid.UUID) error
	delete(id uuid.UUID) error
	getFromPost(postId uuid.UUID) ([]Comment, error)
	get(id uuid.UUID) (Comment, error)
}

type commentRepository struct{}

func (i *commentRepository) create(comment Comment) (uuid.UUID, error) {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return uuid.UUID{}, err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	var id uuid.UUID
	err = tx.QueryRow(
		context.Background(),
		"INSERT INTO comments (user_id, post_id, description) VALUES ($1, $2, $3) RETURNING id",
		comment.UserID, comment.PostID, comment.Description,
	).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert comment: %w", err)
	}

	return id, nil
}

func (i *commentRepository) update(comment Comment, id uuid.UUID) error {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	_, err = tx.Exec(
		context.Background(),
		"UPDATE comments SET description = $1 WHERE id = $1",
		comment.Description, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

func (i *commentRepository) like(id uuid.UUID) error {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	_, err = tx.Exec(
		context.Background(),
		"UPDATE comments SET like_count = like_count + 1 WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

func (i *commentRepository) unlike(id uuid.UUID) error {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	_, err = tx.Exec(
		context.Background(),
		"UPDATE comments SET like_count = like_count - 1 WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

func (i *commentRepository) delete(id uuid.UUID) error {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	_, err = tx.Exec(context.Background(), "DELETE FROM comments WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (i *commentRepository) getFromPost(postId uuid.UUID) ([]Comment, error) {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	rows, err := tx.Query(
		context.Background(),
		"SELECT id, user_id, description, like_count, created_at FROM comments WHERE post_id = $1 ORDER BY like_count DESC",
		postId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to select comments: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.UserID, &comment.Description, &comment.LikeCount, &comment.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return comments, nil
}

func (i *commentRepository) get(id uuid.UUID) (Comment, error) {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return Comment{}, err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return Comment{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	var comment Comment
	err = tx.QueryRow(
		context.Background(),
		"SELECT id, user_id, post_id FROM comments WHERE id = $1",
		id).Scan(&comment.ID, &comment.UserID, &comment.PostID)
	if err != nil {
		return Comment{}, fmt.Errorf("failed to scan comment: %w", err)
	}

	return comment, nil
}
