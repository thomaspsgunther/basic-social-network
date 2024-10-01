package comments

import (
	"context"
	"fmt"
	"time"
	database "y_net/internal/database/postgres"

	"github.com/google/uuid"
)

type Comments struct {
	CommentList []*Comment
}

type Comment struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	PostID      uuid.UUID `json:"post_id"`
	Description string    `json:"description"`
	LikeCount   int       `json:"like_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func (comment *Comment) Create() error {
	if comment.Description == "" {
		return fmt.Errorf("comment text must not be empty")
	}

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
		"INSERT INTO comments (user_id, post_id, description) VALUES ($1, $2, $3)",
		comment.UserID, comment.PostID, comment.Description,
	)
	if err != nil {
		return fmt.Errorf("failed to insert comment: %w", err)
	}

	return nil
}

func Update(comment Comment, id uuid.UUID) error {
	comment.ID = id

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
		"UPDATE comments SET description = $1 WHERE id = $2",
		comment.Description, comment.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

func Delete(id uuid.UUID) error {
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

func GetFromPost(postId uuid.UUID) ([]Comment, error) {
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

func Get(id uuid.UUID) (Comment, error) {
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
