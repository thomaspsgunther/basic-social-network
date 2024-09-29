package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"y_net/internal/auth"
	"y_net/internal/logger"
	"y_net/internal/models/posts"
)

type PostsResource struct{}

func (rs PostsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", rs.Create) // POST /api/v1/posts - Create a new post
	r.Get("/", rs.List)    // GET /api/v1/posts?limit=10&cursor=base64string - Read a list of posts using pagination

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", rs.Get)       // GET /api/v1/posts/{id} - Read a single post by: id
		r.Post("/", rs.Update)   // POST /api/v1/posts/{id} - Update a single post by: id
		r.Delete("/", rs.Delete) // DELETE /api/v1/posts/{id} - Delete a single post by: id
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/{user_id}", rs.GetFromUser) // GET /api/v1/posts/user/{user_id}?limit=10&cursor=base64string - Read a list of posts by: user_id using pagination
	})

	return r
}

// Request Handler - POST /api/v1/posts - Create a new post
func (rs PostsResource) Create(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: post %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var post posts.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	err = post.Create()
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Request Handler - GET /api/v1/posts?limit=10&cursor=base64string - Read a list of posts using pagination
func (rs PostsResource) List(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: get %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid posts limit", http.StatusBadRequest)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	lastCreatedAt, lastId, err := decodeCursor(cursor)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid posts limit", http.StatusBadRequest)
		return
	}

	posts, err := posts.GetPosts(limit, lastCreatedAt, lastId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(posts)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Request Handler - GET /api/v1/posts/{id} - Read a single post by: id
func (rs PostsResource) Get(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: get %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	postId, err := uuid.Parse(id)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	post, err := posts.GetPost(postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(post)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Request Handler - POST /api/v1/users/{id} - Update a single post by: id
func (rs PostsResource) Update(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: post %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	postId, err := uuid.Parse(id)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	ogPost, err := posts.GetPost(postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if authUser.ID != ogPost.UserID {
		err := fmt.Errorf("forbidden post update attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var post posts.Post
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	err = posts.Update(post, postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Request Handler - DELETE /api/v1/posts/{id} - Delete a single post by: id
func (rs PostsResource) Delete(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: delete %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	postId, err := uuid.Parse(id)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	ogPost, err := posts.GetPost(postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if authUser.ID != ogPost.UserID {
		err := fmt.Errorf("forbidden post delete attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = posts.Delete(postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Request Handler - GET /api/v1/posts/user/{user_id}?limit=10&cursor=base64string - Read a list of posts by: user_id using pagination
func (rs PostsResource) GetFromUser(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: get %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "user_id")
	userId, err := uuid.Parse(id)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid posts limit", http.StatusBadRequest)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	lastCreatedAt, lastId, err := decodeCursor(cursor)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid posts limit", http.StatusBadRequest)
		return
	}

	posts, err := posts.GetFromUser(userId, limit, lastCreatedAt, lastId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(posts)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func decodeCursor(encodedCursor string) (time.Time, uuid.UUID, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return time.Time{}, uuid.UUID{}, err
	}

	arrStr := strings.Split(string(byt), ",")
	if len(arrStr) != 2 {
		return time.Time{}, uuid.UUID{}, fmt.Errorf("invalid posts cursor")
	}

	lastCreatedAt, err := time.Parse(time.RFC3339Nano, arrStr[0])
	if err != nil {
		return time.Time{}, uuid.UUID{}, fmt.Errorf("invalid posts lastCreatedAt")
	}

	lastId, err := uuid.Parse(arrStr[1])
	if err != nil {
		return time.Time{}, uuid.UUID{}, fmt.Errorf("invalid posts lastId")
	}

	return lastCreatedAt, lastId, nil
}
