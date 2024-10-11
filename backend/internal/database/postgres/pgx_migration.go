package database

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"y-net/internal/logger"
	"y-net/internal/utils"
)

type sqlMigration struct {
	file  string
	query string
}

func PgxMigration() error {
	conn, err := Postgres.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	// Check previous migrations
	var existingMigrations []string

	logger.ServerLogger.Info("checking for migrations")

	rows, err := conn.Query(context.Background(), "SELECT file FROM migrations")
	if err != nil {
		return fmt.Errorf("error querying migrations table: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var file string
		if err := rows.Scan(&file); err != nil {
			return err
		}
		existingMigrations = append(existingMigrations, file)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error reading rows: %w", err)
	}

	// Get outstanding migrations
	migrations, err := oustandingMigrations(existingMigrations)
	if err != nil {
		return err
	}

	logger.ServerLogger.Info(fmt.Sprintf("number of migrations: %v", len(migrations)))

	// Run migrations in a transaction
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		HandleTransaction(tx, err)
	}()

	for _, migration := range migrations {
		logger.ServerLogger.Info("executing migration: " + migration.file)

		queries := strings.Split(migration.query, ";")
		for _, query := range queries {
			query = strings.TrimSpace(query)
			if query == "" {
				continue
			}

			if _, err := tx.Exec(context.Background(), query); err != nil {
				return fmt.Errorf("failed to execute query %s: %w", query, err)
			}
		}

		if _, err := tx.Exec(context.Background(), "INSERT INTO migrations (file) VALUES ($1)", migration.file); err != nil {
			return fmt.Errorf("failed to insert migration record: %w", err)
		}
	}

	return nil
}

func oustandingMigrations(migrations []string) ([]sqlMigration, error) {
	migrationsFolder, err := utils.FindFolder("migrations")
	if err != nil {
		return nil, err
	}

	dir, err := os.ReadDir(migrationsFolder)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, fileDir := range dir {
		files = append(files, fileDir.Name())
	}

	test := func(s string) bool {
		return strings.HasSuffix(s, ".sql") && !slices.Contains(migrations, s)
	}
	filesFilter := utils.Filter(files, test)

	var sql []sqlMigration
	for _, file := range filesFilter {
		query, err := os.ReadFile(fmt.Sprintf("%s/%s", migrationsFolder, file))
		if err != nil {
			return nil, err
		}

		sql = append(sql, sqlMigration{file: file, query: string(query)})
	}

	return sql, nil
}
