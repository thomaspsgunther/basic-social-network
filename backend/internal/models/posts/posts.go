package posts

import (
	"context"
	"fmt"
	"time"
	database "y_net/internal/database/postgres"

	"github.com/google/uuid"
)

type Posts struct {
	PostList []*Post
}

type Post struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userId"`
	Image        string    `json:"image"`
	Description  *string   `json:"description"`
	LikeCount    int       `json:"likeCount"`
	CommentCount int       `json:"commentCount"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (post *Post) Create() (uuid.UUID, error) {
	if post.Image == "" {
		return uuid.UUID{}, fmt.Errorf("post image must not be empty")
	}

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
		"INSERT INTO posts (user_id, image, description) VALUES ($1, $2, $3) RETURNING id",
		post.UserID, post.Image, post.Description,
	).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to insert post: %w", err)
	}

	return id, nil
}

func GetPosts(limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]Post, error) {
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

	query := `
		SELECT id, user_id, image, description, like_count, comment_count, created_at 
		FROM posts 
		WHERE (created_at < $1 OR (created_at = $1 AND id < $2))
		ORDER BY created_at DESC, id DESC 
		LIMIT $3
	`
	rows, err := tx.Query(context.Background(), query, lastCreatedAt, lastId, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to select posts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Image, &post.Description, &post.LikeCount, &post.CommentCount, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return posts, nil
}

func GetPost(id uuid.UUID) (Post, error) {
	conn, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return Post{}, err
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return Post{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	var post Post
	err = tx.QueryRow(
		context.Background(),
		"SELECT id, user_id, image, description, like_count FROM posts WHERE id = $1",
		id).Scan(&post.ID, &post.UserID, &post.Image, &post.Description, &post.LikeCount)
	if err != nil {
		return Post{}, fmt.Errorf("failed to scan post: %w", err)
	}

	return post, nil
}

func Update(post Post, id uuid.UUID) error {
	post.ID = id

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
		"UPDATE posts SET description = $1 WHERE id = $2",
		post.Description, post.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
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

	_, err = tx.Exec(context.Background(), "DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func GetFromUser(userId uuid.UUID, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]Post, error) {
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

	query := `
        SELECT id, user_id, image, description, like_count, comment_count, created_at 
        FROM posts 
        WHERE user_id = $1 
        AND (created_at < $2 OR (created_at = $2 AND id < $3))
        ORDER BY created_at DESC, id DESC 
        LIMIT $4
    `

	rows, err := tx.Query(context.Background(), query, userId, lastCreatedAt, lastId, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to select posts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Image, &post.Description, &post.LikeCount, &post.CommentCount, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return posts, nil
}
