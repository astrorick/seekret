package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/astrorick/seekret/internal/api"
)

// implement username checks here
func isUsernameValid(username string) bool {
	// TODO: add strict username characters checks
	if len(username) < 4 || len(username) > 32 {
		return false
	}

	return true
}

// implement password checks here
func isPasswordValid(password string) bool {
	// TODO: add strict password characters checks
	if len(password) < 8 || len(password) > 64 {
		return false
	}

	return true
}

// TODO: doc
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

		// prepare srp for salt and verifier generation
		salt, err := srv.SRPParams.NewSalt()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "internal srp salt generation error",
			})
			return
		}
		verifier, err := srv.SRPParams.GetVerifier(salt, newUser.Username, newUser.Password)
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
		if err := srv.Database.CreateUser(newUser.Username, salt, verifier); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(api.ServerErrorResponse{
				Outcome: "failed",
				Reason:  "internal database write error",
			})
			return
		}

		// generate a jwt for the newly created user
		signedJWTString, err := srv.newSignedJWTString(newUser.Username, 24*time.Hour)
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
