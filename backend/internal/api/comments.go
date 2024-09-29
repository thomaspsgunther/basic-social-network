package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"y_net/internal/auth"
	"y_net/internal/logger"
	"y_net/internal/models/comments"
	"y_net/internal/models/posts"
)

type CommentsResource struct{}

func (rs CommentsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", rs.Create)              // POST /api/v1/comments - Create a new comment
	r.Get("/{post_id}", rs.GetFromPost) // GET /api/v1/comments/{post_id} - Read a list of comments by: post_id

	r.Route("/{id}", func(r chi.Router) {
		r.Post("/", rs.Update)   // POST /api/v1/comments/{id} - Update a single comment by: id
		r.Delete("/", rs.Delete) // DELETE /api/v1/comments/{id} - Delete a single comment by: id
	})

	return r
}

// Request Handler - POST /api/v1/comments - Create a new comment
func (rs CommentsResource) Create(w http.ResponseWriter, r *http.Request) {
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

		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// Request Handler - GET /api/v1/comments/{post_id} - Read a list of comments by: post_id
func (rs CommentsResource) GetFromPost(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: get %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "post_id")
	postId, err := uuid.Parse(id)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	comments, err := comments.GetFromPost(postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(comments)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Request Handler - POST /api/v1/comments/{id} - Update a single comment by: id
func (rs CommentsResource) Update(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: post %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var comment comments.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := chi.URLParam(r, "id")
	commentId, err := uuid.Parse(id)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId := comment.UserID

	if authUser.ID != userId {
		err := fmt.Errorf("forbidden comment update attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = comments.Update(comment, commentId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Request Handler - DELETE /api/v1/comments/{id} - Delete a single comment by: id
func (rs CommentsResource) Delete(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: delete %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	commentId, err := uuid.Parse(id)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	comment, err := comments.Get(commentId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post, err := posts.Get(comment.PostID)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if authUser.ID != comment.UserID || authUser.ID != post.UserID {
		err := fmt.Errorf("forbidden comment delete attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = comments.Delete(commentId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
