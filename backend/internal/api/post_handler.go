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

	"y-net/internal/auth"
	"y-net/internal/logger"
	"y-net/internal/services/posts"
	"y-net/internal/services/shared"
)

type PostHandler struct {
	Usecase posts.IPostUsecase
}

func (h PostHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.CreatePost)                            // POST /api/v1/posts - Create a new post
	r.Get("/", h.ListPosts)                              // GET /api/v1/posts?limit=10&cursor=base64string - Read a list of posts using pagination
	r.Post("/likes/{user_id}_{post_id}", h.Like)         // POST /api/v1/posts/likes/{user_id}_{post_id} - Like a post by: id
	r.Delete("/likes/{user_id}_{post_id}", h.Unlike)     // DELETE /api/v1/posts/likes/{user_id}_{post_id} - Unlike a post by: id
	r.Get("/check/{user_id}_{post_id}", h.UserLikedPost) // GET /api/v1/posts/checkliked/{user_id}_{post_id} - Check if a user has liked a post by: id
	r.Get("/likes/{id}", h.GetLikes)                     // GET /api/v1/posts/likes/{id} - Read a list of users who liked a post by: post_id

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.GetPost)       // GET /api/v1/posts/{id} - Read a single post by: id
		r.Put("/", h.UpdatePost)    // PUT /api/v1/posts/{id} - Update a single post by: id
		r.Delete("/", h.DeletePost) // DELETE /api/v1/posts/{id} - Delete a single post by: id
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
// @Param       body body shared.Post true "Post Object"
// @Success     200 {object} shared.Post
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

	var post shared.Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if authUser.ID != post.User.ID {
		err := fmt.Errorf("forbidden post create attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	id, err := h.Usecase.Create(r.Context(), post)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(shared.Post{ID: id})
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
// @Param       cursor query string false "cursor for pagination" Format(byte)
// @Success     200 {object} shared.Posts
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

	var posts []shared.Post
	cursor := r.URL.Query().Get("cursor")
	if cursor != "" {
		lastCreatedAt, lastId, err := decodeCursor(cursor)
		if err != nil {
			logger.ServerLogger.Error(err.Error())

			http.Error(w, "invalid posts cursor", http.StatusBadRequest)
			return
		}

		posts, err = h.Usecase.GetPosts(r.Context(), limit, lastCreatedAt, lastId)
		if err != nil {
			logger.ServerLogger.Error(err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		posts, err = h.Usecase.GetPosts(r.Context(), limit, time.Time{}, uuid.Nil)
		if err != nil {
			logger.ServerLogger.Error(err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

// Like         godoc
// @Summary     Like a post by: id
// @Description Like a post by: id
// @Tags        posts
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       user_id path string true "User ID" Format(uuid)
// @Param       post_id path string true "Post ID" Format(uuid)
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     403
// @Failure     500
// @Router      /posts/likes/{user_id}_{post_id} [post]
func (h PostHandler) Like(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: post %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userId, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if authUser.ID != userId {
		err := fmt.Errorf("forbidden like attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	postId, err := uuid.Parse(chi.URLParam(r, "post_id"))
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	err = h.Usecase.Like(r.Context(), userId, postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Unlike       godoc
// @Summary     Unlike a post by: id
// @Description Unlike a post by: id
// @Tags        posts
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       user_id path string true "User ID" Format(uuid)
// @Param       post_id path string true "Post ID" Format(uuid)
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     403
// @Failure     500
// @Router      /posts/likes/{user_id}_{post_id} [delete]
func (h PostHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: delete %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userId, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if authUser.ID != userId {
		err := fmt.Errorf("forbidden unlike attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	postId, err := uuid.Parse(chi.URLParam(r, "post_id"))
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	err = h.Usecase.Unlike(r.Context(), userId, postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UserLikedPost godoc
// @Summary      Check if a user has liked a post by: id
// @Description  Check if a user has liked a post by: id
// @Tags         posts
// @Produce      json
// @Param        Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param        user_id path string true "User ID" Format(uuid)
// @Param        post_id path string true "Post ID" Format(uuid)
// @Success      200 {object} posts.LikedJson
// @Failure      400
// @Failure      401
// @Failure      500
// @Router       /posts/checkliked/{user_id}_{post_id} [get]
func (h PostHandler) UserLikedPost(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: get %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userId, err := uuid.Parse(chi.URLParam(r, "user_id"))
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	postId, err := uuid.Parse(chi.URLParam(r, "post_id"))
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	isLiked, err := h.Usecase.UserLikedPost(r.Context(), userId, postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(posts.LikedJson{Liked: isLiked})
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetLikes         godoc
// @Summary         Read a list of users who liked a post by: post_id
// @Description     Read a list of users who liked a post by: post_id
// @Tags            posts
// @Produce         json
// @Param           Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param           post_id path string true "Post ID" Format(uuid)
// @Success         200 {object} shared.Users
// @Failure         400
// @Failure         401
// @Failure         500
// @Router          /posts/likes/{id} [get]
func (h PostHandler) GetLikes(w http.ResponseWriter, r *http.Request) {
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

	likes, err := h.Usecase.GetLikes(r.Context(), postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(likes)
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
// @Success     200 {object} shared.Post
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

	post, err := h.Usecase.GetPost(r.Context(), postId)
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
// @Param       body body shared.Post true "Post Object"
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     403
// @Failure     500
// @Router      /posts/{id} [put]
func (h PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: put %s", r.URL))

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

	ogPost, err := h.Usecase.GetPost(r.Context(), postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if authUser.ID != ogPost.User.ID {
		err := fmt.Errorf("forbidden post update attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var post shared.Post
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.Usecase.Update(r.Context(), post, postId)
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

	ogPost, err := h.Usecase.GetPost(r.Context(), postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if authUser.ID != ogPost.User.ID {
		err := fmt.Errorf("forbidden post delete attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = h.Usecase.Delete(r.Context(), postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func decodeCursor(encodedCursor string) (time.Time, uuid.UUID, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}

	arrStr := strings.Split(string(byt), ",")
	if len(arrStr) != 2 {
		return time.Time{}, uuid.Nil, fmt.Errorf("invalid posts cursor")
	}

	lastCreatedAt, err := time.Parse(time.RFC3339Nano, arrStr[0])
	if err != nil {
		return time.Time{}, uuid.Nil, fmt.Errorf("invalid posts lastCreatedAt")
	}

	lastId, err := uuid.Parse(arrStr[1])
	if err != nil {
		return time.Time{}, uuid.Nil, fmt.Errorf("invalid posts lastId")
	}

	return lastCreatedAt, lastId, nil
}
