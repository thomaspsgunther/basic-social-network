package users

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	database "y-net/internal/database/postgres"
	"y-net/internal/services/shared"
)

type iUserRepository interface {
	create(ctx context.Context, user shared.User) (uuid.UUID, error)
	get(ctx context.Context, idList []uuid.UUID) ([]shared.User, error)
	getBySearch(ctx context.Context, searchStr string) ([]shared.User, error)
	getPostsFromUser(ctx context.Context, userId uuid.UUID, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]shared.Post, error)
	update(ctx context.Context, user shared.User, id uuid.UUID) error
	delete(ctx context.Context, id uuid.UUID) error
	follow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error
	unfollow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error
	userFollowsUser(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) (bool, error)
	getFollowers(ctx context.Context, iId uuid.UUID) ([]shared.User, error)
	getFollowed(ctx context.Context, id uuid.UUID) ([]shared.User, error)
}

type userRepositoryImpl struct{}

func (r *userRepositoryImpl) create(ctx context.Context, user shared.User) (uuid.UUID, error) {
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
		"INSERT INTO users (username, password, email, full_name, description, avatar) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		user.Username, user.Password, user.Email, user.FullName, user.Description, user.Avatar,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *userRepositoryImpl) get(ctx context.Context, idList []uuid.UUID) ([]shared.User, error) {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	var placeholders []string
	args := make([]interface{}, len(idList))
	for i, id := range idList {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args[i] = id
	}

	placeholdersClause := strings.Join(placeholders, ", ")
	query := fmt.Sprintf("SELECT id, username, full_name, avatar, post_count, follower_count, followed_count FROM users WHERE id IN (%s)", placeholdersClause)

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select users: %w", err)
	}
	defer rows.Close()

	var users []shared.User
	for rows.Next() {
		var user shared.User
		if err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Avatar, &user.PostCount, &user.FollowerCount, &user.FollowedCount); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return users, nil
}

func (r *userRepositoryImpl) getBySearch(ctx context.Context, searchStr string) ([]shared.User, error) {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	query := "SELECT id, username, full_name, avatar FROM users WHERE username ILIKE $1 OR full_name ILIKE $1"

	searchPattern := "%" + searchStr + "%"

	rows, err := tx.Query(ctx, query, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to select users: %w", err)
	}
	defer rows.Close()

	var users []shared.User
	for rows.Next() {
		var user shared.User
		if err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Avatar); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return users, nil
}

func (r *userRepositoryImpl) getPostsFromUser(ctx context.Context, userId uuid.UUID, limit int, lastCreatedAt time.Time, lastId uuid.UUID) ([]shared.Post, error) {
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
			SELECT id, image, created_at
			FROM posts
			WHERE user_id = $1
			ORDER BY created_at DESC, id DESC
			LIMIT $2
		`
		args = append(args, userId, limit)
	} else {
		query = `
			SELECT id, image, created_at
			FROM posts
			WHERE user_id = $1
			AND (created_at < $2 OR (created_at = $2 AND id < $3))
			ORDER BY created_at DESC, id DESC
			LIMIT $4
		`
		args = append(args, userId, lastCreatedAt, lastId, limit)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []shared.Post{}, nil
		}

		return nil, fmt.Errorf("failed to select posts: %w", err)
	}
	defer rows.Close()

	var posts []shared.Post
	for rows.Next() {
		var post shared.Post
		if err := rows.Scan(&post.ID, &post.Image, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return posts, nil
}

func (r *userRepositoryImpl) update(ctx context.Context, user shared.User, id uuid.UUID) error {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	if user.Password != "" {
		_, err = tx.Exec(
			ctx,
			"UPDATE users SET username = $1, password = $2, email = $3, full_name = $4, description = $5, avatar = $6 WHERE id = $7",
			user.Username, user.Password, user.Email, user.FullName, user.Description, user.Avatar, id,
		)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		_, err = tx.Exec(
			ctx,
			"UPDATE users SET username = $1, email = $2, full_name = $3, description = $4, avatar = $5 WHERE id = $6",
			user.Username, user.Email, user.FullName, user.Description, user.Avatar, id,
		)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	return nil
}

func (r *userRepositoryImpl) delete(ctx context.Context, id uuid.UUID) error {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	_, err = tx.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *userRepositoryImpl) follow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	_, err = tx.Exec(
		ctx,
		"INSERT INTO followers (follower_id, followed_id) VALUES ($1, $2)",
		followerId, followedId,
	)
	if err != nil {
		return fmt.Errorf("failed to insert follower: %w", err)
	}

	return nil
}

func (r *userRepositoryImpl) unfollow(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) error {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	_, err = tx.Exec(ctx, "DELETE FROM followers WHERE follower_id = $1 AND followed_id = $2", followerId, followedId)
	if err != nil {
		return fmt.Errorf("failed to delete follower: %w", err)
	}

	return nil
}

func (r *userRepositoryImpl) userFollowsUser(ctx context.Context, followerId uuid.UUID, followedId uuid.UUID) (bool, error) {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	var exists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM likes WHERE follower_id = $1 AND followed_id = $2)", followerId, followedId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user follows user: %w", err)
	}

	return exists, nil
}

func (r *userRepositoryImpl) getFollowers(ctx context.Context, id uuid.UUID) ([]shared.User, error) {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	rows, err := tx.Query(ctx, "SELECT u.id, u.username, u.full_name, u.avatar FROM followers f JOIN users u ON f.follower_id = u.id WHERE f.followed_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to select followers: %w", err)
	}
	defer rows.Close()

	var userFollowers []shared.User
	for rows.Next() {
		var user shared.User
		if err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Avatar); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		userFollowers = append(userFollowers, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return userFollowers, nil
}

func (r *userRepositoryImpl) getFollowed(ctx context.Context, id uuid.UUID) ([]shared.User, error) {
	tx, err := database.Postgres.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		database.HandleTransaction(ctx, tx, err)
	}()

	rows, err := tx.Query(ctx, "SELECT u.id, u.username, u.full_name, u.avatar FROM followers f JOIN users u ON f.followed_id = u.id WHERE f.follower_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to select followed: %w", err)
	}
	defer rows.Close()

	var userFollowed []shared.User
	for rows.Next() {
		var user shared.User
		if err := rows.Scan(&user.ID, &user.Username, &user.FullName, &user.Avatar); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		userFollowed = append(userFollowed, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return userFollowed, nil
}
