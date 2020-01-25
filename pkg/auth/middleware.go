package auth

import ()

import (
	"github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

func getProtectedHandler() {
	jwtMiddleware.Handler(CertHandler)
}
