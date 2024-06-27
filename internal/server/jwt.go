package server

import (
	"github.com/golang-jwt/jwt/v4"
)

type JWTParams struct {
	JWTHashFcn jwt.SigningMethod
	JWTKey     []byte
}

func NewParams(hashFcn string, key string) (*JWTParams, error) {
	return nil, nil
}

/*func (srv *Server) newSignedJWTString(username string, duration time.Duration) (string, error) {
	// create JWT token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(duration).Unix(),
	})

	// sign the token with the server super secret key
	signedJWTString, err := jwtToken.SignedString(srv.Key)
	if err != nil {
		return "", err
	}

	return signedJWTString, nil
}*/
