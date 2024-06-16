package api

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
