package seekret

import "math/big"

type SRPClient struct {
	username string
	salt     []byte
	a        *big.Int
	A        *big.Int
	B        *big.Int
}

type SRPServer struct {
}
