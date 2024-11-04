package posts

import (
	"context"
	"fmt"
	"time"
	database "y-net/internal/database/postgres"
	"y-net/internal/services/shared"

	"github.com/google/uuid"
)

type iPostRepository interface {
	create(ctx context.Context, post shared.Post) (uuid.UUID, error)
	getPosts(ctx context.Context, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]shared.Post, error)
	getPost(ctx context.Context, id uuid.UUID) (shared.Post, error)
	update(ctx context.Context, post shared.Post, id uuid.UUID) error
	delete(ctx context.Context, id uuid.UUID) error
	like(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error
	unlike(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error
	userLikedPost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error)
	getLikes(ctx context.Context, id uuid.UUID) ([]shared.User, error)
}

type postRepositoryImpl struct{}

func (r *postRepositoryImpl) create(ctx context.Context, post shared.Post) (uuid.UUID, error) {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
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
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	var query string
	var args []interface{}

	if lastCreatedAt.IsZero() && lastId == uuid.Nil {
		query = `
			SELECT p.id, p.user_id, u.username, u.avatar, p.image, p.description, p.like_count, p.comment_count, p.created_at
			FROM posts p
			INNER JOIN users u ON p.user_id = u.id
			ORDER BY p.created_at DESC, p.id DESC
			LIMIT $1
		`
		args = append(args, limit)
	} else {
		query = `
			SELECT p.id, p.user_id, u.username, u.avatar, p.image, p.description, p.like_count, p.comment_count, p.created_at
			FROM posts p
			INNER JOIN users u ON p.user_id = u.id
			WHERE (p.created_at < $1 OR (p.created_at = $1 AND p.id < $2))
			ORDER BY p.created_at DESC, p.id DESC
			LIMIT $3
		`
		args = append(args, lastCreatedAt, lastId, limit)
	}

	rows, err := tx.Query(ctx, query, args...)
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
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return shared.Post{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	query := `
		SELECT p.id, p.user_id, u.username, u.avatar, p.image, p.description, p.like_count, p.comment_count, p.created_at
		FROM posts p
		INNER JOIN users u ON p.user_id = u.id
		WHERE p.id = $1
	`

	var post shared.Post
	err = tx.QueryRow(
		ctx,
		query,
		id,
	).Scan(&post.ID, &post.User.ID, &post.User.Username, &post.User.Avatar, &post.Image, &post.Description, &post.LikeCount, &post.CommentCount, &post.CreatedAt)
	if err != nil {
		return shared.Post{}, fmt.Errorf("failed to scan post: %w", err)
	}

	return post, nil
}

func (r *postRepositoryImpl) update(ctx context.Context, post shared.Post, id uuid.UUID) error {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
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
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	_, err = tx.Exec(ctx, "DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (i *postRepositoryImpl) like(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	_, err = tx.Exec(
		ctx,
		"INSERT INTO likes (user_id, post_id) VALUES ($1, $2)",
		userId, postId,
	)
	if err != nil {
		return fmt.Errorf("failed to insert like: %w", err)
	}

	return nil
}

func (i *postRepositoryImpl) unlike(ctx context.Context, userId uuid.UUID, postId uuid.UUID) error {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	_, err = tx.Exec(ctx, "DELETE FROM likes WHERE user_id = $1 AND post_id = $2", userId, postId)
	if err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}

	return nil
}

func (i *postRepositoryImpl) userLikedPost(ctx context.Context, userId uuid.UUID, postId uuid.UUID) (bool, error) {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND post_id = $2)", userId, postId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user liked post: %w", err)
	}

	return exists, nil
}

func (i *postRepositoryImpl) getLikes(ctx context.Context, id uuid.UUID) ([]shared.User, error) {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	rows, err := tx.Query(ctx, "SELECT u.id, u.username, u.full_name, u.avatar FROM likes l JOIN users u ON l.user_id = u.id WHERE l.post_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to select likes: %w", err)
	}
	defer rows.Close()

	var userLikes []shared.User
	for rows.Next() {
		var user shared.User
		if err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Avatar); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		userLikes = append(userLikes, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return userLikes, nil
}
