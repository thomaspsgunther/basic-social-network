package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"y_net/internal/auth"
	"y_net/internal/logger"
	"y_net/internal/models/followers"
)

type FollowersResource struct{}

func (rs FollowersResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/{follower_id}_{followed_id}", rs.Follow)     // POST /api/v1/followers/{follower_id}_{followed_id} - Follow a user by: id
	r.Delete("/{follower_id}_{followed_id}", rs.Unfollow) // DELETE /api/v1/followers/{follower_id}_{followed_id} - Unfollow a user by: id
	r.Get("/{user_id}", rs.GetFollowers)                  // GET /api/v1/followers/{id} - Read a list of who follows a user by: user_id

	r.Route("/followed", func(r chi.Router) {
		r.Get("/{user_id}", rs.GetFollowed) // GET /api/v1/followers/followed/{id} - Read a list of who a user follows by: user_id
	})

	return r
}

// Follow       godoc
// @Summary     Follow a user by: id
// @Description Follow a user by: id
// @Tags        followers
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       follower_id path string true "Follower ID" Format(uuid)
// @Param       followed_id path string true "Followed ID" Format(uuid)
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     403
// @Failure     500
// @Router      /followers/{follower_id}_{followed_id} [post]
func (rs FollowersResource) Follow(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: post %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	followerId, err := uuid.Parse(chi.URLParam(r, "follower_id"))
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if authUser.ID != followerId {
		err := fmt.Errorf("forbidden follow attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	followedId, err := uuid.Parse(chi.URLParam(r, "followed_id"))
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	err = followers.Follow(followerId, followedId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Unfollow     godoc
// @Summary     Unfollow a user by: id
// @Description Unfollow a user by: id
// @Tags        followers
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       follower_id path string true "Follower ID" Format(uuid)
// @Param       followed_id path string true "Followed ID" Format(uuid)
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     403
// @Failure     500
// @Router      /followers/{follower_id}_{followed_id} [delete]
func (rs FollowersResource) Unfollow(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: delete %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	followerId, err := uuid.Parse(chi.URLParam(r, "follower_id"))
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	followedId, err := uuid.Parse(chi.URLParam(r, "followed_id"))
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if authUser.ID != followerId || authUser.ID != followedId {
		err := fmt.Errorf("forbidden unfollow attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = followers.Unfollow(followerId, followedId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetFollowers godoc
// @Summary     Read a list of who follows a user by: user_id
// @Description Read a list of who follows a user by: user_id
// @Tags        followers
// @Produce     json
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       user_id path string true "User ID" Format(uuid)
// @Success     200 {object} users.Users
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /followers/{user_id} [get]
func (rs FollowersResource) GetFollowers(w http.ResponseWriter, r *http.Request) {
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

	followers, err := followers.GetFollowers(userId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(followers)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetFollowed  godoc
// @Summary     Read a list of who a user follows by: user_id
// @Description Read a list of who a user follows by: user_id
// @Tags        followers
// @Produce     json
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       user_id path string true "User ID" Format(uuid)
// @Success     200 {object} users.Users
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /followers/followed/{user_id} [get]
func (rs FollowersResource) GetFollowed(w http.ResponseWriter, r *http.Request) {
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

	followed, err := followers.GetFollowed(userId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(followed)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
