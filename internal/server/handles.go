package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/astrorick/seekret/internal/api"
	"github.com/golang-jwt/jwt/v4"
)

func isUsernameValid(username string) bool {
	// TODO: add strict username characters checks
	if len(username) < 4 || len(username) > 32 {
		return false
	}

	return true
}

func isPasswordValid(password string) bool {
	// TODO: add strict password characters checks
	if len(password) < 8 || len(password) > 64 {
		return false
	}

	return true
}

func (srv *Server) newSignedJWT(username string, duration time.Duration) (string, error) {
	// create JWT token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(duration).Unix(),
	})

	// sign the token with the server super secret key
	signedJWTString, err := jwtToken.SignedString(srv.JWTKey)
	if err != nil {
		return "", err
	}

	return signedJWTString, nil
}

func (srv *Server) CreateUserRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check for the correct method
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "method not allowed",
			})
			return
		}

		// read request body
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "unable to read request body",
			})
			return
		}

		// parse request body
		var newUser api.CreateUserRequest
		if err := json.Unmarshal(reqBody, &newUser); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "unable to parse request content",
			})
			return
		}

		// check parsed data
		if !isUsernameValid(newUser.Username) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "invalid username",
			})
			return
		}
		if !isPasswordValid(newUser.Password) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "invalid password",
			})
			return
		}

		// check if specified username exists
		userExists, err := srv.Database.UserExists(newUser.Username)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "internal database read error",
			})
			return
		}
		if userExists {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "username not available",
			})
			return
		}

		// add new user
		// TODO: SRP salt and verifier
		if err := srv.Database.CreateUser(newUser.Username, newUser.Password, newUser.Password); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "internal database write error",
			})
			return
		}

		// generate a JWT for the newly created user
		signedJWTString, err := srv.newSignedJWT(newUser.Username, 24*time.Hour)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "internal JWT generator error",
			})
			return
		}

		// send feedback
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(api.CreateUserResponse{
			JWT: signedJWTString,
		})
	}
}
