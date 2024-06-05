package api

type LoginClientRequest struct {
	Identity string
	A        string
}

type LoginServerResponse struct {
	Salt string
	B    string
}

type ClientProofRequest struct {
	Identity string
	M1       string
}

type ServerProofResponse struct {
	M2  string
	JWT string
}
