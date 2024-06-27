package srp

import (
	"crypto"
	"crypto/rand"
)

type Params struct {
	SaltSize uint64      // size in bytes
	HashFcn  crypto.Hash // hash function
	/*N        *big.Int    // selected safe prime
	G        *big.Int    // generator modulo N*/
}

// NewSalt generates a new random salt of length specified by the SRP parameters.
func (p *Params) NewSalt() ([]byte, error) {
	salt := make([]byte, p.SaltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	return salt, nil
}

// GetVerifier returns the verifier from the specified arguments and using the hashing function provided in the parameters.
// This function should only be used during the creation process of a new user and not for authentication.
func (p *Params) GetVerifier(salt []byte, username string, password string) ([]byte, error) {
	h1 := p.HashFcn.New()
	h2 := p.HashFcn.New()

	// compute hash of [username][":"][password] as h1
	if _, err := h1.Write([]byte(username + ":" + password)); err != nil {
		return nil, err
	}

	// compute hash of [salt][h1] as h2
	if _, err := h2.Write(append(salt, h1.Sum(nil)...)); err != nil {
		return nil, err
	}

	return h2.Sum(nil), nil
}
