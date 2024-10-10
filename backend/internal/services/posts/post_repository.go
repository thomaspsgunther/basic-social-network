package posts

import (
	"context"
	"fmt"
	"time"
	database "y_net/internal/database/postgres"
	"y_net/internal/services/shared"

	"github.com/google/uuid"
)

type postRepository interface {
	create(ctx context.Context, post shared.Post) (uuid.UUID, error)
	getPosts(ctx context.Context, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]shared.Post, error)
	getPost(ctx context.Context, id uuid.UUID) (shared.Post, error)
	update(ctx context.Context, post shared.Post, id uuid.UUID) error
	delete(ctx context.Context, id uuid.UUID) error
}

type postRepositoryImpl struct{}

func (r *postRepositoryImpl) create(ctx context.Context, post shared.Post) (uuid.UUID, error) {
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
		"INSERT INTO posts (user_id, image, description) VALUES ($1, $2, $3) RETURNING id",
		post.User.ID, post.Image, post.Description,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert post: %w", err)
	}

	return id, nil
}

func (r *postRepositoryImpl) getPosts(ctx context.Context, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]shared.Post, error) {
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
		SELECT p.id, p.user_id, u.username, u.avatar, p.image, p.description, p.like_count, p.comment_count, p.created_at
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		WHERE (created_at < $1 OR (created_at = $1 AND id < $2))
		ORDER BY created_at DESC, id DESC
		LIMIT $3
	`
	rows, err := tx.Query(ctx, query, lastCreatedAt, lastId, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to select posts: %w", err)
	}
	defer rows.Close()

	var posts []shared.Post
	for rows.Next() {
		var post shared.Post
		err := rows.Scan(&post.ID, &post.User.ID, &post.User.Username, &post.User.Avatar, &post.Image, &post.Description, &post.LikeCount, &post.CommentCount, &post.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return posts, nil
}

func (r *postRepositoryImpl) getPost(ctx context.Context, id uuid.UUID) (shared.Post, error) {
	conn, err := database.Postgres.Acquire(ctx)
	if err != nil {
		return shared.Post{}, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return shared.Post{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(tx, err)
	}()

	query := `
		SELECT p.id, p.user_id, u.username, u.avatar, p.image, p.description, p.like_count, p.comment_count, p.created_at
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		WHERE id = $1
	`

	var post shared.Post
	err = tx.QueryRow(
		ctx,
		query,
		id,
	).Scan(&post.ID, &post.User.ID, &post.User.Username, &post.User.Avatar, &post.Image, &post.Description, &post.LikeCount)
	if err != nil {
		return shared.Post{}, fmt.Errorf("failed to scan post: %w", err)
	}

	return post, nil
}

func (r *postRepositoryImpl) update(ctx context.Context, post shared.Post, id uuid.UUID) error {
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
		"UPDATE posts SET description = $1 WHERE id = $2",
		post.Description, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}

func (r *postRepositoryImpl) delete(ctx context.Context, id uuid.UUID) error {
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

	_, err = tx.Exec(ctx, "DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
