package jwt

import (
	"github.com/golang-jwt/jwt/v4"
)

type Params struct {
	SigningFcn jwt.SigningMethod
	Key        []byte
}

func (p *Params) NewSignedJWTString(mapClaims jwt.MapClaims) (string, error) {
	// create JWT token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)

	// sign the token with the server super secret key
	signedJWTString, err := jwtToken.SignedString(p.Key)
	if err != nil {
		return "", err
	}

	return signedJWTString, nil
}
