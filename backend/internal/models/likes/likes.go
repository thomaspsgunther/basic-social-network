package likes

import (
	"context"
	"fmt"
	database "y_net/internal/database/postgres"
	"y_net/internal/models/users"

	"github.com/google/uuid"
)

type Like struct {
	UserID uuid.UUID `json:"user_id"`
	PostID uuid.UUID `json:"post_id"`
}

func LikePost(userId uuid.UUID, postId uuid.UUID) error {
	like := Like{UserID: userId, PostID: postId}

	connection, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer connection.Release()

	tx, err := connection.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO likes (user_id, post_id) VALUES ($1, $2)",
		like.UserID, like.PostID,
	)
	if err != nil {
		return fmt.Errorf("failed to insert like: %w", err)
	}

	return nil
}

func UnlikePost(userId uuid.UUID, postId uuid.UUID) error {
	connection, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer connection.Release()

	tx, err := connection.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	_, err = tx.Exec(context.Background(), "DELETE FROM likes WHERE user_id = $1 AND post_id = $2", userId, postId)
	if err != nil {
		return fmt.Errorf("failed to delete like: %w", err)
	}

	return nil
}

func GetFromPost(postId uuid.UUID) ([]users.User, error) {
	connection, err := database.Postgres.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer connection.Release()

	tx, err := connection.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer database.HandleTransaction(tx, err)

	rows, err := tx.Query(context.Background(), "SELECT u.id, u.username, u.full_name, u.avatar FROM likes l JOIN users u ON l.user_id = u.id WHERE l.post_id = $1", postId)
	if err != nil {
		return nil, fmt.Errorf("failed to select likes: %w", err)
	}
	defer rows.Close()

	var followers []users.User
	for rows.Next() {
		var follower users.User
		if err := rows.Scan(&follower.ID, &follower.Username, &follower.FullName, &follower.Avatar); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		followers = append(followers, follower)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return followers, nil
}
