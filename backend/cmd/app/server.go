package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/zishang520/socket.io/v2/socket"

	"y_net/internal/api"
	"y_net/internal/auth"
	database "y_net/internal/database/postgres"
	"y_net/internal/logger"
	"y_net/internal/socketio"
	"y_net/internal/utils"
)

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
	connectPG, err := utils.GetBoolFromString(os.Getenv("PG_CONN"))
	if err != nil {
		logger.ServerLogger.Info("--------------------------------------------------------------------")
		logger.ServerLogger.Error(err.Error())
	}

	// DB connect
	if connectPG {
		err := database.PostgresConnect()
		if err != nil {
			logger.ServerLogger.Info("--------------------------------------------------------------------")
			logger.ServerLogger.Fatalf("failed to connect PostgreSQL: %v", err)
		}
		defer database.PostgresClose()
	}

	// DB migrations
	if connectPG {
		err := database.PostgresMigration()
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
	router := chi.NewRouter()
	router.Use(auth.Middleware())
	router.Handle("/socket.io/", io.ServeHandler(nil))
	router.HandleFunc("/api/v1/", rootFunc)
	router.Mount("/api/v1/login", api.LoginResource{}.Routes())
	router.Mount("/api/v1/users", api.UsersResource{}.Routes())
	router.Mount("/api/v1/followers", api.FollowersResource{}.Routes())
	router.Mount("/api/v1/posts", api.PostsResource{}.Routes())
	router.Mount("/api/v1/likes", api.LikesResource{}.Routes())
	router.Mount("/api/v1/comments", api.CommentsResource{}.Routes())

	// Start the server api
	host := os.Getenv("HTTP_HOST")
	port := os.Getenv("HTTP_PORT")

	logger.ServerLogger.Info(fmt.Sprintf("server running on http://%s:%s/api/v1/", host, port))
	logger.ServerLogger.Info("--------------------------------------------------------------------")

	go http.ListenAndServe(":"+port, router)

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
