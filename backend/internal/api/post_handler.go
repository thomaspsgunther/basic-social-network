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
	"y_net/internal/services/posts"
)

type PostHandler struct {
	Usecase posts.PostUsecaseI
}

func (h PostHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.CreatePost) // POST /api/v1/posts - Create a new post
	r.Get("/", h.ListPosts)   // GET /api/v1/posts?limit=10&cursor=base64string - Read a list of posts using pagination

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.GetPost)       // GET /api/v1/posts/{id} - Read a single post by: id
		r.Post("/", h.UpdatePost)   // POST /api/v1/posts/{id} - Update a single post by: id
		r.Delete("/", h.DeletePost) // DELETE /api/v1/posts/{id} - Delete a single post by: id
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/{user_id}", h.ListPostsFromUser) // GET /api/v1/posts/user/{user_id}?limit=10&cursor=base64string - Read a list of posts by: user_id using pagination
	})

	return r
}

// CreatePost   godoc
// @Summary     Create a new post
// @Description Create a new post
// @Tags        posts
// @Accept      json
// @Produce     json
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       body body posts.Post true "Post Object"
// @Success     200 {object} posts.Post
// @Failure     400
// @Failure     401
// @Failure     403
// @Failure     500
// @Router      /posts [post]
func (h PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
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

	if authUser.ID != post.UserID {
		err := fmt.Errorf("forbidden post create attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	id, err := h.Usecase.Create(post)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(posts.Post{ID: id})
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// ListPosts    godoc
// @Summary     Read a list of posts using pagination
// @Description Read a list of posts using pagination
// @Tags        posts
// @Produce     json
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       limit query int true "limit of pagination"
// @Param       cursor query string true "cursor for pagination" Format(byte)
// @Success     200 {object} posts.Posts
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /posts [get]
func (h PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
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

	posts, err := h.Usecase.GetPosts(limit, lastCreatedAt, lastId)
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

// GetPost      godoc
// @Summary     Read a single post by: id
// @Description Read a single post by: id
// @Tags        posts
// @Produce     json
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       id path string true "Post ID" Format(uuid)
// @Success     200 {object} posts.Post
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /posts/{id} [get]
func (h PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
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

	post, err := h.Usecase.GetPost(postId)
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

// UpdatePost   godoc
// @Summary     Update a single post by: id
// @Description Update a single post by: id
// @Tags        posts
// @Accept      json
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       id path string true "Post ID" Format(uuid)
// @Param       body body posts.Post true "Post Object"
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     403
// @Failure     500
// @Router      /posts/{id} [post]
func (h PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
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

	ogPost, err := h.Usecase.GetPost(postId)
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

	err = h.Usecase.Update(post, postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeletePost   godoc
// @Summary     Delete a single post by: id
// @Description Delete a single post by: id
// @Tags        posts
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       id path string true "Post ID" Format(uuid)
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     403
// @Failure     500
// @Router      /posts/{id} [delete]
func (h PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
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

	ogPost, err := h.Usecase.GetPost(postId)
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

	err = h.Usecase.Delete(postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ListPostsFromUser godoc
// @Summary          Read a list of posts by: user_id using pagination
// @Description      Read a list of posts by: user_id using pagination
// @Tags             posts
// @Produce          json
// @Param            Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param            user_id path string true "User ID" Format(uuid)
// @Param            limit query int true "limit of pagination"
// @Param            cursor query string true "cursor for pagination" Format(byte)
// @Success          200 {object} posts.Posts
// @Failure          400
// @Failure          401
// @Failure          500
// @Router           /posts/user/{user_id} [get]
func (h PostHandler) ListPostsFromUser(w http.ResponseWriter, r *http.Request) {
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

	posts, err := h.Usecase.GetFromUser(userId, limit, lastCreatedAt, lastId)
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
