package api

type LoginClientRequest struct {
	Username string
	A        string
}

type LoginServerResponse struct {
	Salt string
	B    string
}

type ClientProofRequest struct {
	Username string
	M1       string
}

type ServerProofResponse struct {
	M2  string
	JWT string
}
