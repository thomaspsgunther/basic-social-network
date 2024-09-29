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

	r.Post("/{follower_id}_{followed_id}", rs.Follow) // POST /api/v1/followers/{follower_id}-{followed_id} - Follow a user by: id
	r.Get("/{id}", rs.GetFollowers)                   // GET /api/v1/followers/{id} - Read a list of followers of the user by: id

	r.Route("/unfollow", func(r chi.Router) {
		r.Delete("/{follower_id}_{followed_id}", rs.Unfollow) // DELETE /api/v1/followers/unfollow/{follower_id}-{followed_id} - Unfollow a user by: id
	})

	r.Route("/followed", func(r chi.Router) {
		r.Get("/{id}", rs.GetFollowed) // GET /api/v1/followers/followed/{id} - Read a list of who the user follows by: id
	})

	return r
}

// Request Handler - POST /api/v1/followers/{follower_id}-{followed_id} - Follow a user by: id
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

		http.Error(w, err.Error(), http.StatusInternalServerError)
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

		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// Request Handler - GET /api/v1/followers/{id} - Read a list of followers of the user by: id
func (rs FollowersResource) GetFollowers(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: get %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	userId, err := uuid.Parse(id)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// Request Handler - DELETE /api/v1/followers/unfollow/{follower_id}-{followed_id} - Unfollow a user by: id
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

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	followedId, err := uuid.Parse(chi.URLParam(r, "followed_id"))
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// Request Handler - GET /api/v1/followers/followed/{id} - Read a list of who the user follows by: id
func (rs FollowersResource) GetFollowed(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: get %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	userId, err := uuid.Parse(id)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
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
