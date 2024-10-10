package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/zishang520/socket.io/v2/socket"

	_ "y_net/docs"
	"y_net/internal/api"
	"y_net/internal/auth"
	database "y_net/internal/database/postgres"
	"y_net/internal/logger"
	"y_net/internal/services/comments"
	"y_net/internal/services/posts"
	"y_net/internal/services/users"
	"y_net/internal/socketio"
	"y_net/internal/utils"
)

// @title        Y API
// @version      1.0
// @description  API for a basic social network.
// @license.name MIT
// @license.url  https://opensource.org/license/mit
// @host         localhost:8080
// @BasePath     /api/v1/
func main() {
	// Start logging
	err := logger.InitLogger("daily")
	if err != nil {
		log.Println("--------------------------------------------------------------------")
		log.Fatal(err)
	}
	defer logger.CloseLogger()

	// Load .env file
	err = utils.LoadEnv()
	if err != nil {
		logger.ServerLogger.Info("--------------------------------------------------------------------")
		logger.ServerLogger.Fatal(err)
	}

	// Check whether to connect to PostgreSQL or not
	connectPG, err := strconv.ParseBool(os.Getenv("PG_CONN"))
	if err != nil {
		logger.ServerLogger.Info("--------------------------------------------------------------------")
		logger.ServerLogger.Error(err.Error())
	}

	// DB connect
	if connectPG {
		err := database.PgxConnect()
		if err != nil {
			logger.ServerLogger.Info("--------------------------------------------------------------------")
			logger.ServerLogger.Fatalf("failed to connect PostgreSQL: %v", err)
		}
		defer database.PgxClose()

		err = database.PgxMigration()
		if err != nil {
			logger.ServerLogger.Info("--------------------------------------------------------------------")
			logger.ServerLogger.Fatalf("failed PostgreSQL migrations: %v", err)
		}
	}

	// Socketio setup
	io := socket.NewServer(nil, nil)
	defer io.Close(nil)
	socketio.ConnectionHandler(io)

	// Routes setup
	r := chi.NewRouter()
	r.Use(auth.Middleware())
	r.Handle("/socket.io/", io.ServeHandler(nil))
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	r.HandleFunc("/api/v1/", rootFunc)
	r.Mount("/api/v1/login", api.LoginHandler{Usecase: users.NewUserUsecase()}.Routes())
	r.Mount("/api/v1/users", api.UserHandler{Usecase: users.NewUserUsecase()}.Routes())
	r.Mount("/api/v1/posts", api.PostHandler{Usecase: posts.NewPostUsecase()}.Routes())
	r.Mount("/api/v1/comments", api.CommentHandler{Usecase: comments.NewCommentUsecase()}.Routes())

	// Start the server api
	host := os.Getenv("HTTP_HOST")
	port := os.Getenv("HTTP_PORT")

	logger.ServerLogger.Info(fmt.Sprintf("server running on http://%s:%s/api/v1/", host, port))
	logger.ServerLogger.Info("--------------------------------------------------------------------")

	go http.ListenAndServe(":"+port, r)

	// Shutdown server
	exit := make(chan struct{})
	SignalC := make(chan os.Signal, 1)

	signal.Notify(SignalC, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range SignalC {
			switch s {
			case os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				close(exit)
				return
			}
		}
	}()

	<-exit
	logger.ServerLogger.Info("--------------------------------------------------------------------")
	logger.ServerLogger.Warn("shutting down...")
	os.Exit(0)
}

func rootFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
