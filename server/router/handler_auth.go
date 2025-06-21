package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Thybaau/todolist-app/middleware"
	"golang.org/x/crypto/bcrypt"
)

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Todo-List by Thibault")
	}
}

func (s *server) handleSignUp() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			middleware.NewHTTPError(w, "Cannot decode task body from json", http.StatusBadRequest, err)
			return
		}
		if req.Username == "" || req.Password == "" {
			middleware.NewHTTPError(w, "Username or password cannot be empty", http.StatusForbidden, nil)
			return
		}

		// Check if user already exist
		_, err = s.DB.GetUserInfos(req.Username)
		if err != nil && err != sql.ErrNoRows {
			middleware.NewHTTPError(w, "Error when searching username", http.StatusInternalServerError, err)
			return
		}
		if err == nil {
			message := fmt.Sprintf("User %v already exist", req.Username)
			middleware.NewHTTPError(w, message, http.StatusForbidden, err)
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			middleware.NewHTTPError(w, "Failed to hash password", http.StatusInternalServerError, err)
			return
		}

		// Register user in database
		_, err = s.DB.CreateUser(req.Username, hashedPassword)
		if err != nil {
			middleware.NewHTTPError(w, "Failed to create account", http.StatusInternalServerError, err)
			return
		}

		successMessage := fmt.Sprintf("successfully signed up with username %v", req.Username)
		jsonResp := map[string]string{"message": successMessage}
		middleware.JSONResponse(w, http.StatusOK, jsonResp)
	}
}

func (s *server) handleLogin() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			middleware.NewHTTPError(w, "Cannot decode task body from json", http.StatusBadRequest, err)
			return
		}
		if req.Username == "" || req.Password == "" {
			middleware.NewHTTPError(w, "Username or password cannot be empty", http.StatusForbidden, nil)
			return
		}

		user, err := s.DB.GetUserInfos(req.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				middleware.NewHTTPError(w, "Username not found", http.StatusNotFound, err)
				return
			} else {
				middleware.NewHTTPError(w, "Error when searching username", http.StatusInternalServerError, err)
				return
			}
		}
		err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(req.Password))
		if err != nil {
			middleware.NewHTTPError(w, "Incorrect password", http.StatusForbidden, err)
			return
		}

		// Generate JWT token
		token, err := middleware.GenerateToken(req.Username)
		if err != nil {
			middleware.NewHTTPError(w, "Cannot generate token", http.StatusInternalServerError, err)
			return
		}

		// Response
		successMessage := fmt.Sprintf("successfully logged as %v", req.Username)
		jsonResp := map[string]string{"message": successMessage, "token": token}
		middleware.JSONResponse(w, http.StatusOK, jsonResp)
	}
}
