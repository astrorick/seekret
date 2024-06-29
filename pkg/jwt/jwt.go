package jwt

import (
	"github.com/golang-jwt/jwt/v4"
)

/*
Exporting these vars and types is required to avoid importing the
third-party jwt package in each other file that uses jwt.
*/

var SigningMethodHS256 = jwt.SigningMethodHS256
var SigningMethodHS384 = jwt.SigningMethodHS384
var SigningMethodHS512 = jwt.SigningMethodHS512

type MapClaims = jwt.MapClaims

// Params defines a set of parameters to be used when creating or validating jwts, like signing function and signing key.
type Params struct {
	SigningFcn jwt.SigningMethod
	Key        []byte
}

// NewSignedJWTString returns a string representation of a jwt generated from the specified parameters and with user provided claims.
// It returns an error if the jwt token cannot be signed (if the provided json is broken, for example).
func (p *Params) NewSignedJWTString(mapClaims MapClaims) (string, error) {
	// create jwt token object
	jwtToken := jwt.NewWithClaims(p.SigningFcn, mapClaims)

	// sign the token with the super secret server key
	signedJWTString, err := jwtToken.SignedString(p.Key)
	if err != nil {
		return "", err
	}

	return signedJWTString, nil
}
