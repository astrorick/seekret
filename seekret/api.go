package seekret

import (
	"encoding/json"
	"io"
	"net/http"
)

// this is used to create a new user in the server database
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (srv *Server) CreateUserRequestHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check for the correct method
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// parse request body
		var newUser CreateUserRequest
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "unable to read request body", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(reqBody, &newUser); err != nil {
			http.Error(w, "unable to parse request content", http.StatusBadRequest)
			return
		}

		// check parsed data
		if newUser.Username == "" {
			http.Error(w, "username not valid", http.StatusBadRequest)
			return
		}
		if newUser.Password == "" {
			http.Error(w, "password not valid", http.StatusBadRequest)
			return
		}

		// check if specified username exists
		var usernameCount int
		if err := srv.Database.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", newUser.Username).Scan(&usernameCount); err != nil {
			http.Error(w, "unable to read database", http.StatusInternalServerError)
			return
		}
		if usernameCount > 0 {
			http.Error(w, "username not available", http.StatusConflict)
			return
		}

		// add new user
		if _, err := srv.Database.Exec("INSERT INTO users (username, salt, verifier) VALUES (?, ?, ?)", newUser.Username, newUser.Password, newUser.Password); err != nil {
			http.Error(w, "unable to write to database", http.StatusInternalServerError)
			return
		}

		// send feedback
		w.WriteHeader(http.StatusCreated)
	}
}

// this is the request the client makes to the server when initiating SRP authentication
type LoginClientRequest struct {
	Username string
	A        string
}

// this is the server response to the first client request
type LoginServerResponse struct {
	SessionID string
	Salt      string
	B         string
}

// this is the request the client makes when providing it's SRP proof
type ClientProofRequest struct {
	SessionID string
	M1        string
}

// this is the server response to the client's proof request
type ServerProofResponse struct {
	M2  string
	JWT string
}
