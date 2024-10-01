package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"y_net/internal/auth"
	"y_net/internal/logger"
	"y_net/internal/models/users"
	"y_net/pkg/jwt"
)

type UsersResource struct{}

func (rs UsersResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", rs.CreateUser)       // POST /api/v1/users - Create a new user
	r.Get("/{id_list}", rs.GetUsers) // GET /api/v1/users/{id_list} - Read a list of users by: id_list

	r.Route("/{id}", func(r chi.Router) {
		r.Post("/", rs.UpdateUser)   // POST /api/v1/users/{id} - Update a single user by: id
		r.Delete("/", rs.DeleteUser) // DELETE /api/v1/users/{id} - Delete a single user by: id
	})

	r.Route("/search", func(r chi.Router) {
		r.Get("/{search_term}", rs.SearchUsers) // GET /api/v1/users/search/{search_term} - Read a list of users by: search_term
	})

	return r
}

// CreateUser   godoc
// @Summary     Create a new user
// @Description Create a new user
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       body body users.User true "User Object"
// @Success     200 {object} users.TokenJson
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /users [post]
func (rs UsersResource) CreateUser(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: post %s", r.URL))

	var user users.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	err = user.Create()
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := users.GetUserIdByUsername(user.Username)
	if err != nil {
		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	tokenStr, err := jwt.GenerateToken(id)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(users.TokenJson{Token: tokenStr})
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetUsers     godoc
// @Summary     Read a list of users by: id_list
// @Description Read a list of users by: id_list
// @Tags        users
// @Produce     json
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       id_list path string true "User ID List"
// @Success     200 {object} users.Users
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /users/{id_list} [get]
func (rs UsersResource) GetUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := users.Get(idList)
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
// @Success     200 {object} users.Users
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /users/search/{search_term} [get]
func (rs UsersResource) SearchUsers(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: get %s", r.URL))

	authUser := auth.ForContext(r.Context())
	if authUser == nil {
		err := fmt.Errorf("access denied")

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	searchStr := chi.URLParam(r, "search_term")

	users, err := users.GetUsersBySearch(searchStr)
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

// UpdateUser   godoc
// @Summary     Update a single user by: id
// @Description Update a single user by: id
// @Tags        users
// @Accept      json
// @Param       Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param       id path string true "User ID" Format(uuid)
// @Param       body body users.User true "User Object"
// @Success     200
// @Failure     400
// @Failure     401
// @Failure     403
// @Failure     500
// @Router      /users/{id} [post]
func (rs UsersResource) UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	var user users.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	err = users.Update(user, userId)
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
func (rs UsersResource) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	err = users.Delete(userId)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
