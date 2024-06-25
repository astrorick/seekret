package srp

import "math/big"

type SRPClient struct {
	SRPParams *SRPParams
	Username  string
	Password  string
	a         *big.Int
	A         *big.Int
}

func (p *SRPParams) NewClient(username string, password string, salt []byte) (*SRPClient, error) {
	// TODO: compute 'a' and 'A'

	return &SRPClient{
		SRPParams: p,
		Username:  username,
		Password:  password,
		a:         &big.Int{},
		A:         &big.Int{},
	}, nil
}
