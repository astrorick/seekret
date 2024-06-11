package seekret

import "math/big"

type SRPParams struct {
	// TODO
	params string
}

type SRPClient struct {
	Params   *SRPParams
	Username string
	Password string
	a        *big.Int
	A        *big.Int
}

func NewClient(params *SRPParams, username string, password string, salt []byte) (*SRPClient, error) {
	// TODO: compute 'a' and 'A'

	return &SRPClient{
		Params:   params,
		Username: username,
		Password: password,
		a:        &big.Int{},
		A:        &big.Int{},
	}, nil
}

type SRPServer struct {
}
