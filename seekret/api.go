package seekret

import (
	"encoding/json"
	"io"
	"net/http"
)

// this is the request the client makes when creating a new user in the server database
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// this is the request the client makes to the server when initiating SRP authentication
type LoginClientRequest struct {
	Username string `json:"username"`
	A        string `json:"A"`
}

// this is the server response to the first client request
type LoginServerResponse struct {
	SessionID string `json:"sessionID"`
	Salt      string `json:"salt"`
	B         string `json:"B"`
}

// this is the request the client makes when providing it's SRP proof
type ClientProofRequest struct {
	SessionID string `json:"sessionID"`
	M1        string `json:"M1"`
}

// this is the server response to the client's proof request
type ServerProofResponse struct {
	M2  string `json:"M2"`
	JWT string `json:"JWT"`
}

// this is the server response to standard client requests
type OutcomeResponse struct {
	Outcome string `json:"outcome"`
	Reason  string `json:"reason"`
}

func (srv *Server) CreateUserRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check for the correct method
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(OutcomeResponse{
				Outcome: "failed",
				Reason:  "method not allowed",
			})
			return
		}

		// parse request body
		var newUser CreateUserRequest
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(OutcomeResponse{
				Outcome: "failed",
				Reason:  "unable to read request body",
			})
			return
		}
		if err := json.Unmarshal(reqBody, &newUser); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(OutcomeResponse{
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
			json.NewEncoder(w).Encode(OutcomeResponse{
				Outcome: "failed",
				Reason:  "invalid username",
			})
			return
		}
		if newUser.Password == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(OutcomeResponse{
				Outcome: "failed",
				Reason:  "invalid password",
			})
			return
		}

		// check if specified username exists
		var checkUsername int
		if err := srv.Database.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", newUser.Username).Scan(&checkUsername); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(OutcomeResponse{
				Outcome: "failed",
				Reason:  "internal database read error",
			})
			return
		}
		if checkUsername > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(OutcomeResponse{
				Outcome: "failed",
				Reason:  "username not available",
			})
			return
		}

		// add new user
		// TODO: SRP salt and verifier
		if _, err := srv.Database.Exec("INSERT INTO users (username, salt, verifier) VALUES (?, ?, ?)", newUser.Username, newUser.Password, newUser.Password); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(OutcomeResponse{
				Outcome: "failed",
				Reason:  "internal database write error",
			})
			return
		}

		// send feedback
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(OutcomeResponse{
			Outcome: "success",
			Reason:  "new user created",
		})
	}
}
