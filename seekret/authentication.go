package seekret

type LoginClientRequest struct {
	Username string
	A        string
}

type LoginServerResponse struct {
	SessionID string
	Salt      string
	B         string
}

type ClientProofRequest struct {
	SessionID string
	M1        string
}

type ServerProofResponse struct {
	M2  string
	JWT string
}
