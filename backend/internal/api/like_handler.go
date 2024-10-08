package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"y_net/internal/auth"
	"y_net/internal/logger"
	"y_net/internal/services/likes"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type LikeHandler struct {
	Usecase likes.LikeUsecaseI
}

func (h LikeHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/{user_id}_{post_id}", h.LikePost)           // POST /api/v1/likes/{user_id}_{post_id} - Like a post by: id
	r.Delete("/{user_id}_{post_id}", h.UnlikePost)       // DELETE /api/v1/likes/{user_id}_{post_id} - Unlike a post by: id
	r.Get("/check/{user_id}_{post_id}", h.UserLikedPost) // GET /api/v1/likes/{user_id}_{post_id} - Check if a user has liked a post by: id
	r.Get("/{post_id}", h.GetLikesFromPost)              // GET /api/v1/likes/{post_id} - Read a list of users who liked a post by: post_id

	return r
}

// LikePost     godoc
// @Summary     Like a post by: id
// @Description Like a post by: id
// @Tags        likes
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       user_id path string true "User ID" Format(uuid)
// @Param       post_id path string true "Post ID" Format(uuid)
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     403
// @Failure     500
// @Router      /likes/{user_id}_{post_id} [post]
func (h LikeHandler) LikePost(w http.ResponseWriter, r *http.Request) {
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

	err = h.Usecase.LikePost(userId, postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UnlikePost   godoc
// @Summary     Unlike a post by: id
// @Description Unlike a post by: id
// @Tags        likes
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       user_id path string true "User ID" Format(uuid)
// @Param       post_id path string true "Post ID" Format(uuid)
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     403
// @Failure     500
// @Router      /likes/{user_id}_{post_id} [delete]
func (h LikeHandler) UnlikePost(w http.ResponseWriter, r *http.Request) {
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

	err = h.Usecase.UnlikePost(userId, postId)
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
// @Tags         likes
// @Produce      json
// @Param        Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param        user_id path string true "User ID" Format(uuid)
// @Param        post_id path string true "Post ID" Format(uuid)
// @Success      200 {object} likes.LikedJson
// @Failure      400
// @Failure      401
// @Failure      500
// @Router       /likes/check/{user_id}_{post_id} [get]
func (h LikeHandler) UserLikedPost(w http.ResponseWriter, r *http.Request) {
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

	isLiked, err := h.Usecase.UserLikedPost(userId, postId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(likes.LikedJson{Liked: isLiked})
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetLikesFromPost godoc
// @Summary         Read a list of users who liked a post by: post_id
// @Description     Read a list of users who liked a post by: post_id
// @Tags            likes
// @Produce         json
// @Param           Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param           post_id path string true "Post ID" Format(uuid)
// @Success         200 {object} users.Users
// @Failure         400
// @Failure         401
// @Failure         500
// @Router          /likes/{post_id} [get]
func (h LikeHandler) GetLikesFromPost(w http.ResponseWriter, r *http.Request) {
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

	likes, err := h.Usecase.GetFromPost(postId)
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
