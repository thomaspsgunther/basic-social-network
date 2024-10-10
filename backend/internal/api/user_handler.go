package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"y_net/internal/auth"
	"y_net/internal/logger"
	"y_net/internal/services/shared"
	"y_net/internal/services/users"
)

type UserHandler struct {
	Usecase users.UserUsecase
}

func (h UserHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/{id_list}", h.GetUsers)                // GET /api/v1/users/{id_list} - Read a list of users by: id_list
	r.Get("/search/{search_term}", h.SearchUsers)  // GET /api/v1/users/search/{search_term} - Read a list of users by: search_term
	r.Get("/posts/{user_id}", h.ListPostsFromUser) // GET /api/v1/users/posts/{user_id}?limit=10&cursor=base64string - Read a list of posts by: user_id using pagination

	r.Route("/{id}", func(r chi.Router) {
		r.Put("/", h.UpdateUser)    // PUT /api/v1/users/{id} - Update a single user by: id
		r.Delete("/", h.DeleteUser) // DELETE /api/v1/users/{id} - Delete a single user by: id
	})

	return r
}

// GetUsers     godoc
// @Summary     Read a list of users by: id_list
// @Description Read a list of users by: id_list
// @Tags        users
// @Produce     json
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       id_list path string true "User ID List"
// @Success     200 {object} shared.Users
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /users/{id_list} [get]
func (h UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: get %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	idListOG := chi.URLParam(r, "id_list")
	idListStr := strings.Split(idListOG, ",")

	var idList []uuid.UUID
	for _, id := range idListStr {
		userId, err := uuid.Parse(id)
		if err != nil {
			logger.ServerLogger.Error(err.Error())

			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		idList = append(idList, userId)
	}

	users, err := h.Usecase.Get(r.Context(), idList)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// SearchUsers  godoc
// @Summary     Read a list of users by: search_term
// @Description Read a list of users by: search_term
// @Tags        users
// @Produce     json
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       search_term path string true "User Search Term"
// @Success     200 {object} shared.Users
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /users/search/{search_term} [get]
func (h UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: get %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	searchStr := chi.URLParam(r, "search_term")

	users, err := h.Usecase.GetBySearch(r.Context(), searchStr)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// ListPostsFromUser godoc
// @Summary          Read a list of posts by: user_id using pagination
// @Description      Read a list of posts by: user_id using pagination
// @Tags             users
// @Produce          json
// @Param            Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param            user_id path string true "User ID" Format(uuid)
// @Param            limit query int true "limit of pagination"
// @Param            cursor query string true "cursor for pagination" Format(byte)
// @Success          200 {object} shared.Posts
// @Failure          400
// @Failure          401
// @Failure          500
// @Router           /users/posts/{user_id} [get]
func (h UserHandler) ListPostsFromUser(w http.ResponseWriter, r *http.Request) {
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

	posts, err := h.Usecase.GetPostsFromUser(r.Context(), userId, limit, lastCreatedAt, lastId)
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

// UpdateUser   godoc
// @Summary     Update a single user by: id
// @Description Update a single user by: id
// @Tags        users
// @Accept      json
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       id path string true "User ID" Format(uuid)
// @Param       body body shared.User true "User Object"
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     403
// @Failure     500
// @Router      /users/{id} [put]
func (h UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: post %s", r.URL))

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

		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if authUser.ID != userId {
		err := fmt.Errorf("forbidden user update attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var user shared.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.Usecase.Update(r.Context(), user, userId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteUser   godoc
// @Summary     Delete a single user by: id
// @Description Delete a single user by: id
// @Tags        users
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       id path string true "User ID" Format(uuid)
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     403
// @Failure     500
// @Router      /users/{id} [delete]
func (h UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: delete %s", r.URL))

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

		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	if authUser.ID != userId {
		err := fmt.Errorf("forbidden user delete attempt from user: %v", authUser.ID)

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = h.Usecase.Delete(r.Context(), userId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
