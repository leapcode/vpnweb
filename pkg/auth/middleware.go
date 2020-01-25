package auth

import (
	"0xacab.org/leap/vpnweb/pkg/web"
	"github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
)

const anonAuth string = "anon"
const sipAuth string = "sip"

/* FIXME -- get this from configuration variables */
var jwtSecret = []byte("somethingverysecret")

func Authenticator(auth string) {
}

func RestrictedMiddleware(auth string, ch web.CertHandler) http.Handler {

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	switch auth {
	case anonAuth:
		return http.HandlerFunc(ch.CertResponder)
	case sipAuth:
		return jwtMiddleware.Handler(http.HandlerFunc(ch.CertResponder))
	default:
		log.Fatal("Unknown auth module: '", auth, "'. Should be one of: ", anonAuth, ", ", sipAuth, ".")
	}
	// should not get here
	return nil
}
