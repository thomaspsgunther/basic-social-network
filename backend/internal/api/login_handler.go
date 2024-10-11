package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"y-net/internal/logger"
	"y-net/internal/services/shared"
	"y-net/internal/services/users"
	"y-net/pkg/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type LoginHandler struct {
	Usecase users.IUserUsecase
}

func (h LoginHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/register", h.CreateUser)       // POST /api/v1/login/register - Create a new user
	r.Post("/", h.Login)                    // POST /api/v1/login - Login user
	r.Post("/refreshtoken", h.RefreshToken) // POST /api/v1/login/refreshtoken - Refresh user token

	return r
}

// CreateUser   godoc
// @Summary     Create a new user
// @Description Create a new user
// @Tags        login
// @Accept      json
// @Produce     json
// @Param       body body shared.User true "User Object"
// @Success     200 {object} shared.TokenJson
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /login/register [post]
func (h LoginHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: post %s", r.URL))

	var user shared.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	id, err := h.Usecase.Create(r.Context(), user)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenStr, err := jwt.GenerateToken(id)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(shared.TokenJson{Token: tokenStr})
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Login        godoc
// @Summary     Login user
// @Description Login user
// @Tags        login
// @Accept      json
// @Produce     json
// @Param       body body shared.User true "User Object"
// @Success     200 {object} shared.TokenJson
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /login [post]
func (h LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: post %s", r.URL))

	var user shared.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	correct, err := users.Authenticate(r.Context(), user)
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

	id, err := users.GetUserIdByUsername(r.Context(), user.Username)
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

	response, err := json.Marshal(shared.TokenJson{Token: tokenStr})
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
// @Param       body body shared.TokenJson true "Token Object"
// @Success     200 {object} shared.TokenJson
// @Failure     400
// @Failure     401
// @Failure     500
// @Router      /login/refreshtoken [post]
func (h LoginHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	logger.ServerLogger.Info(fmt.Sprintf("new request: post %s", r.URL))

	var token shared.TokenJson
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

	response, err := json.Marshal(shared.TokenJson{Token: tokenStr})
	if err != nil {
		logger.ServerLogger.Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
