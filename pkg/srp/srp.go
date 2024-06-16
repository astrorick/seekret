package seekret

import (
	"crypto"
	"math/big"
)

type SRP struct {
	hashFcn *crypto.Hash
}

type SRPClient struct {
	SRP      *SRP
	Username string
	Password string
	a        *big.Int
	A        *big.Int
}

func NewClient(srp *SRP, username string, password string, salt []byte) (*SRPClient, error) {
	// TODO: compute 'a' and 'A'

	return &SRPClient{
		SRP:      srp,
		Username: username,
		Password: password,
		a:        &big.Int{},
		A:        &big.Int{},
	}, nil
}

type SRPServer struct {
}
