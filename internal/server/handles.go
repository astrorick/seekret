package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/astrorick/seekret/internal/api"
)

func (srv *Server) CreateUserRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check for the correct method
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(api.OutcomeResponse{
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
			json.NewEncoder(w).Encode(api.OutcomeResponse{
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
			json.NewEncoder(w).Encode(api.OutcomeResponse{
				Outcome: "failed",
				Reason:  "unable to parse request content",
			})
			return
		}

		// check parsed data
		// TODO: add more checks for username and password
		if newUser.Username == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.OutcomeResponse{
				Outcome: "failed",
				Reason:  "invalid username",
			})
			return
		}
		if newUser.Password == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(api.OutcomeResponse{
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
			json.NewEncoder(w).Encode(api.OutcomeResponse{
				Outcome: "failed",
				Reason:  "internal database read error",
			})
			return
		}
		if userExists {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(api.OutcomeResponse{
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
			json.NewEncoder(w).Encode(api.OutcomeResponse{
				Outcome: "failed",
				Reason:  "internal database write error",
			})
			return
		}

		// send feedback
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(api.OutcomeResponse{
			Outcome: "success",
			Reason:  "new user created",
		})
	}
}
