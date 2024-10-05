package database

import (
	"context"
	"fmt"
	"y_net/internal/logger"

	"github.com/jackc/pgx/v5"
)

// HandleTransaction commits a DB transaction if the queries were successful or rolls back if not
func HandleTransaction(tx pgx.Tx, err error) {
	if r := recover(); r != nil {
		tx.Rollback(context.Background())

		logger.ServerLogger.Error(fmt.Sprintf("panic occurred during transaction: %v", r))
	} else if err != nil {
		tx.Rollback(context.Background())
	} else {
		if err = tx.Commit(context.Background()); err != nil {
			logger.ServerLogger.Error("failed to commit transaction: " + err.Error())
		}
	}
}
