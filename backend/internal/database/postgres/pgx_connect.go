package database

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
	"y-net/internal/logger"

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
	defaultMaxConns, err := strconv.ParseInt(os.Getenv("PG_CONN_MAX"), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to create a pgxpool config: %w", err)
	}
	defaultMinConns, err := strconv.ParseInt(os.Getenv("PG_CONN_MIN"), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to create a pgxpool config: %w", err)
	}
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
		return nil, fmt.Errorf("failed to create a pgxpool config: %w", err)
	}

	dbConfig.MaxConns = int32(defaultMaxConns)
	dbConfig.MinConns = int32(defaultMinConns)
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	return dbConfig, nil
}
