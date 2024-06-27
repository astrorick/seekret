package srp

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"hash"
)

type SRPParams struct {
	SaltSize uint64      // size in bytes
	HashFcn  crypto.Hash // hash function
	/*N        *big.Int    // selected safe prime
	G        *big.Int    // generator modulo N*/
}

func NewParams(salSize uint64, hashFcn crypto.Hash) (*SRPParams, error) {
	return nil, nil
}

// NewSalt generates a new salt of length specified by the SRP parameters.
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
	var h1, h2 hash.Hash

	// define hashing function
	switch p.HashFcn {
	case crypto.SHA256:
		h1 = sha256.New()
		h2 = sha256.New()
	default:
		return nil, fmt.Errorf("unsupported hash function")
	}

	// compute hash of [username][":"][password] as h1
	h1.Write([]byte(username + ":" + password))

	// compute hash of [salt][h1] as h2
	h2.Write(append(salt, h1.Sum(nil)...))

	return h2.Sum(nil), nil
}
