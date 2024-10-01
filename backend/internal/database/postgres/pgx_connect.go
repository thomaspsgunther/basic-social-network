package database

import (
	"context"
	"fmt"
	"os"
	"time"
	"y_net/internal/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Postgres *pgxpool.Pool

func PgxConnect() error {
	pgxPoolConfig, err := postgresConfig()
	if err != nil {
		return err
	}
	Postgres, err = pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil {
		return err
	}

	logger.ServerLogger.Info(fmt.Sprintf("connecting to db at %s, port %s", os.Getenv("PG_CONN_HOST"), os.Getenv("PG_CONN_PORT")))

	return nil
}

func PgxClose() {
	Postgres.Close()
}

func postgresConfig() (*pgxpool.Config, error) {
	const defaultMaxConns = int32(100)
	const defaultMinConns = int32(30)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	databaseConn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("PG_CONN_USER"),
		os.Getenv("PG_CONN_PWD"),
		os.Getenv("PG_CONN_HOST"),
		os.Getenv("PG_CONN_PORT"),
		os.Getenv("PG_CONN_DBNAME"),
	)

	dbConfig, err := pgxpool.ParseConfig(databaseConn)
	if err != nil {
		return nil, fmt.Errorf("failed to create a pgxpool config, error: %w", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
	}

	return dbConfig, nil
}
