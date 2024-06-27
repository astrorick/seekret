package srp

import (
	"crypto"
	"crypto/rand"
)

type SRPParams struct {
	SaltSize uint64      // size in bytes
	HashFcn  crypto.Hash // hash function
	/*N        *big.Int    // selected safe prime
	G        *big.Int    // generator modulo N*/
}

// NewSalt generates a new random salt of length specified by the SRP parameters.
func (p *SRPParams) NewSalt() ([]byte, error) {
	salt := make([]byte, p.SaltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	return salt, nil
}

// GetVerifier returns the verifier from the specified arguments and using the hashing function provided in the parameters.
// This function should only be used during the creation process of a new user and not for authentication.
func (p *SRPParams) GetVerifier(salt []byte, username string, password string) ([]byte, error) {
	h1 := p.HashFcn.New()
	h2 := p.HashFcn.New()

	// compute hash of [username][":"][password] as h1
	h1.Write([]byte(username + ":" + password))

	// compute hash of [salt][h1] as h2
	h2.Write(append(salt, h1.Sum(nil)...))

	return h2.Sum(nil), nil
}
