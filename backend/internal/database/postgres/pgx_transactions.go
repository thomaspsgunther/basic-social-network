package database

import (
	"context"
	"fmt"
	"y-net/internal/logger"

	"github.com/jackc/pgx/v5"
)

// HandleTransaction commits a DB transaction if the queries were successful or rolls back if not
func HandleTransaction(ctx context.Context, tx pgx.Tx, err error) {
	if r := recover(); r != nil {
		tx.Rollback(ctx)

		logger.ServerLogger.Error(fmt.Sprintf("panic occurred during transaction: %v", r))
	} else if err != nil {
		tx.Rollback(ctx)
	} else {
		if err = tx.Commit(ctx); err != nil {
			logger.ServerLogger.Error("failed to commit transaction: " + err.Error())
		}
	}
}
