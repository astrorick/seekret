package srp

import (
	"crypto"
	"math/big"
)

type SRPParams struct {
	NLenghtBits uint64
	HashFcn     crypto.Hash
	G           *big.Int
	N           *big.Int
}

func New(nLengthBits uint64, hashFcn crypto.Hash) (*SRPParams, error) {
	var srpParams *SRPParams

	srpParams.N = knownPrimes[nLengthBits]
	srpParams.HashFcn = hashFcn

	return srpParams, nil
}

type SRPClient struct {
	SRPParams *SRPParams
	Username  string
	Password  string
	a         *big.Int
	A         *big.Int
}

func NewClient(srpParams *SRPParams, username string, password string, salt []byte) (*SRPClient, error) {
	// TODO: compute 'a' and 'A'

	return &SRPClient{
		SRPParams: srpParams,
		Username:  username,
		Password:  password,
		a:         &big.Int{},
		A:         &big.Int{},
	}, nil
}

type SRPServer struct {
}
