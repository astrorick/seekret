package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// Exporting these vars is required to avoid importing the third-party jwt package in each other file that uses jwt.
var (
	SigningMethodHS256 = jwt.SigningMethodHS256
	SigningMethodHS384 = jwt.SigningMethodHS384
	SigningMethodHS512 = jwt.SigningMethodHS512
)

// Exporting these types is required to avoid importing the third-party jwt package in each other file that uses jwt.
type (
	MapClaims = jwt.MapClaims
	Token     = jwt.Token
)

// Params defines a set of parameters to be used when creating or validating jwts, like signing function and signing key.
type Params struct {
	SigningFcn jwt.SigningMethod
	SecretKey  []byte
}

// NewSignedJWTString returns a string representation of a jwt generated from the specified parameters and with user provided claims.
// It returns an error if the jwt token cannot be signed (if the provided json is broken, for example).
func (p *Params) NewSignedJWTString(mapClaims *MapClaims) (string, error) {
	// create jwt token object
	jwtToken := jwt.NewWithClaims(p.SigningFcn, mapClaims)

	// sign the token with the super secret server key
	signedJWTString, err := jwtToken.SignedString(p.SecretKey)
	if err != nil {
		return "", err
	}

	return signedJWTString, nil
}

func (p *Params) ParseJWTString(tokenString string) (*Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check token signing method
		if token.Method.Alg() != p.SigningFcn.Alg() {
			return nil, fmt.Errorf("unexpected jwt signing method")
		}
		return p.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (p *Params) GetClaims(token *Token) (*MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &claims, nil
}
