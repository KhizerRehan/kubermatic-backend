package handler

import (
	"encoding/base64"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
)

func jwtMiddleware(key string) *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			decoded, err := base64.URLEncoding.DecodeString(key)
			if err != nil {
				return nil, err
			}
			return decoded, nil
		},
		Debug: bool(glog.V(6)),
	})
}
