package socketio

import (
	"fmt"
	"y_net/internal/logger"

	"github.com/zishang520/socket.io/v2/socket"
)

func ErrorHandler(client *socket.Socket) {
	client.On("connect_error", func(data ...any) {
		logger.ServerLogger.Error(fmt.Sprint(data))
	})
	client.On("connect_failed", func(data ...any) {
		logger.ServerLogger.Error(fmt.Sprint(data))
	})
}
