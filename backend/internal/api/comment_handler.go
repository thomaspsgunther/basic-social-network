package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"y-net/internal/auth"
	"y-net/internal/logger"
	"y-net/internal/services/comments"
)

type CommentHandler struct {
	Usecase comments.ICommentUsecase
}

func (h CommentHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.CreateComment)               // POST /api/v1/comments - Create a new comment
	r.Get("/{post_id}", h.GetCommentsFromPost) // GET /api/v1/comments/{post_id} - Read a list of comments by: post_id
	r.Post("/like/{id}", h.LikeComment)        // POST /api/v1/comments/like/{id} - Like a single comment by: id
	r.Post("/unlike/{id}", h.UnlikeComment)    // POST /api/v1/comments/unlike/{id} - Unlike a single comment by: id

	r.Route("/{id}", func(r chi.Router) {
		r.Put("/", h.UpdateComment)    // PUT /api/v1/comments/{id} - Update a single comment by: id
		r.Delete("/", h.DeleteComment) // DELETE /api/v1/comments/{id} - Delete a single comment by: id
	})

	return r
}

// CreateComment godoc
// @Summary      Create a new comment
// @Description  Create a new comment
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param        body body comments.Comment true "Comment Object"
// @Success      200 {object} comments.Comment
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      500
// @Router       /comments [post]
func (h CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
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

		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if authUser.ID != comment.User.ID {
		err := fmt.Errorf("forbidden comment create attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	id, err := h.Usecase.Create(r.Context(), comment)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(comments.Comment{ID: id})
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetCommentsFromPost godoc
// @Summary            Read a list of comments by: post_id
// @Description        Read a list of comments by: post_id
// @Tags               comments
// @Produce            json
// @Param              Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param              post_id path string true "Post ID" Format(uuid)
// @Success            200 {object} comments.Comments
// @Failure            401
// @Failure            500
// @Router             /comments/{post_id} [get]
func (h CommentHandler) GetCommentsFromPost(w http.ResponseWriter, r *http.Request) {
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

		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	comments, err := h.Usecase.GetFromPost(r.Context(), postId)
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

// UpdateComment godoc
// @Summary      Update a single comment by: id
// @Description  Update a single comment by: id
// @Tags         comments
// @Accept       json
// @Param        Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param        id path string true "Comment ID" Format(uuid)
// @Param        body body comments.Comment true "Comment Object"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      500
// @Router       /comments/{id} [put]
func (h CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: put %s", r.URL))

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

		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	commentId, err := uuid.Parse(id)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	userId := comment.User.ID

	if authUser.ID != userId {
		err := fmt.Errorf("forbidden comment update attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = h.Usecase.Update(r.Context(), comment, commentId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// LikeComment   godoc
// @Summary      Like a single comment by: id
// @Description  Like a single comment by: id
// @Tags         comments
// @Param        Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param        id path string true "Comment ID" Format(uuid)
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      500
// @Router       /comments/like/{id} [post]
func (h CommentHandler) LikeComment(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: post %s", r.URL))

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

		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	err = h.Usecase.Like(r.Context(), commentId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UnlikeComment godoc
// @Summary      Unlike a single comment by: id
// @Description  Unlike a single comment by: id
// @Tags         comments
// @Param        Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param        id path string true "Comment ID" Format(uuid)
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      500
// @Router       /comments/unlike/{id} [post]
func (h CommentHandler) UnlikeComment(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: post %s", r.URL))

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

		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	err = h.Usecase.Unlike(r.Context(), commentId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteComment godoc
// @Summary      Delete a single comment by: id
// @Description  Delete a single comment by: id
// @Tags         comments
// @Param        Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param        id path string true "Comment ID" Format(uuid)
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      403
// @Failure      500
// @Router       /comments/{id} [delete]
func (h CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
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

		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	comment, err := h.Usecase.Get(r.Context(), commentId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if authUser.ID != comment.User.ID {
		err := fmt.Errorf("forbidden comment delete attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = h.Usecase.Delete(r.Context(), commentId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
