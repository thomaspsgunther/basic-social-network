package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"y_net/internal/logger"
	"y_net/internal/models/users"
	"y_net/pkg/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type LoginResource struct{}

func (rs LoginResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", rs.Login)                    // POST /api/v1/login - Login user
	r.Post("/refreshtoken", rs.RefreshToken) // POST /api/v1/login/refreshtoken - Refresh user token

	return r
}

// Login        godoc
// @Summary     Login user
// @Description Login user
// @Tags        login
// @Accept      json
// @Produce     json
// @Param       body body users.User true "User Object"
// @Success     200 {object} users.TokenJson
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /login [post]
func (rs LoginResource) Login(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: post %s", r.URL))

	var user users.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	correct, err := user.Authenticate()
	if err != nil {
		if err == pgx.ErrNoRows {
			err = &users.WrongUsernameOrPasswordError{}

			logger.ServerLogger.Warn(err.Error())

			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !correct {
		err = &users.WrongUsernameOrPasswordError{}

		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	id, err := users.GetUserIdByUsername(user.Username)
	if err != nil {
		logger.ServerLogger.Warn(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// RefreshToken godoc
// @Summary     Refresh user token
// @Description Refresh user token
// @Tags        login
// @Accept      json
// @Produce     json
// @Param       body body users.TokenJson true "Token Object"
// @Success     200 {object} users.TokenJson
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /login/refreshtoken [post]
func (rs LoginResource) RefreshToken(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: post %s", r.URL))

	var token users.TokenJson
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := jwt.ParseToken(token.Token)
	if err != nil {
		err := fmt.Errorf("access denied")

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
