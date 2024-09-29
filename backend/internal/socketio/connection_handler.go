package socketio

import (
	"fmt"

	"github.com/zishang520/socket.io/v2/socket"

	"y_net/internal/logger"
	"y_net/pkg/jwt"
)

func ConnectionHandler(io *socket.Server) {
	io.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)

		logger.ServerLogger.Info(fmt.Sprintf("connected: %v", client.Id()))

		tokenStr, _ := client.Request().Headers().Get("Authorization")
		if tokenStr == "" {
			logger.ServerLogger.Warn("access denied")

			client.Disconnect(true)
		}
		_, err := jwt.ParseToken(tokenStr)
		if err != nil {
			logger.ServerLogger.Warn("access denied")

			client.Disconnect(true)
		}

		ErrorHandler(client)

		client.On("disconnect", func(data ...any) {
			logger.ServerLogger.Info(fmt.Sprintf("closed %v", data))
		})
	})
}
