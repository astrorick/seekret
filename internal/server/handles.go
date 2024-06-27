package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/astrorick/seekret/internal/api"

	gojwt "github.com/golang-jwt/jwt/v4"
)

// implement username checks here
func isValidUsername(username string) bool {
	// TODO: add strict username characters checks
	if len(username) < 4 || len(username) > 32 {
		return false
	}

	return true
}

// implement password checks here
func isValidPassword(password string) bool {
	// TODO: add strict password characters checks
	if len(password) < 8 || len(password) > 64 {
		return false
	}

	return true
}

// CreateUserRequestHandler handles client request for the creation of a new user in the database.
// An error response is returned to the client if the request body cannot be parsed, if the username/password is not valid or if the user already exists.
// If the request is successful, a new JWT will be returned in the response body to be used by the client.
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
		var newUserRequest api.CreateUserRequest
		if err := json.Unmarshal(reqBody, &newUserRequest); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "unable to parse request content",
			})
			return
		}

		// check parsed data
		if !isValidUsername(newUserRequest.Username) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "invalid username",
			})
			return
		}
		if !isValidPassword(newUserRequest.Password) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "invalid password",
			})
			return
		}

		// check if specified username exists
		userExists, err := srv.Database.UserExists(newUserRequest.Username)
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

		// prepare srp for salt and verifier generation
		newUserSalt, err := srv.SRPParams.NewSalt()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "internal srp salt generation error",
			})
			return
		}
		newUserVerifier, err := srv.SRPParams.GetVerifier(newUserSalt, newUserRequest.Username, newUserRequest.Password)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "internal srp verifier generation error",
			})
			return
		}

		// add new user
		if err := srv.Database.CreateUser(newUserRequest.Username, newUserSalt, newUserVerifier); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "internal database write error",
			})
			return
		}

		// generate a jwt for the newly created user
		signedJWTString, err := srv.JWTParams.NewSignedJWTString(gojwt.MapClaims{
			"username": newUserRequest.Username,
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
		})
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
