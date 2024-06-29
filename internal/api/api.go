package api

/* USER CREATION & DELETION */

// CreateUserRequest is the request the client makes when creating a new user in the server database.
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateUserResponse is the server response when a new user is created.
type CreateUserResponse struct {
	JWT string `json:"JWT"`
}

// DeleteUserRequest is the request the client makes when deleting an existing user from server database.
type DeleteUserRequest struct {
	Username string `json:"username"`
}

/* USER AUTHENTICATION */

// AuthClientRequest is the request the client makes to the server when initiating SRP authentication.
type AuthClientRequest struct {
	Username string `json:"username"`
	A        string `json:"A"`
}

// AuthServerResponse is the server response to the first client request.
type AuthServerResponse struct {
	SessionID string `json:"sessionID"`
	Salt      string `json:"salt"`
	B         string `json:"B"`
}

// ClientProofRequest is the followup request the client makes when providing it's SRP proof.
type ClientProofRequest struct {
	SessionID string `json:"sessionID"`
	M1        string `json:"M1"`
}

// ServerProofResponse is the server response to the followup client's proof request.
type ServerProofResponse struct {
	M2  string `json:"M2"`
	JWT string `json:"JWT"`
}

/* ERRORS */

// ServerErrorResponse is a server error response to failed client requests.
type ServerErrorResponse struct {
	Outcome string `json:"outcome"`
	Reason  string `json:"reason"`
}
