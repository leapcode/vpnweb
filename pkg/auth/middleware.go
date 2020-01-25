package auth

import (
	"0xacab.org/leap/vpnweb/pkg/auth/sip2"
	"0xacab.org/leap/vpnweb/pkg/config"
	"0xacab.org/leap/vpnweb/pkg/web"
	"github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
)

const anonAuth string = "anon"
const sipAuth string = "sip"

/* FIXME -- get this from configuration variables */

var jwtSigningSecret = []byte("thesingingkey")

func bailOnBadAuthModule(module string) {
	log.Fatal("Unknown auth module: '", module, "'. Should be one of: ", anonAuth, ", ", sipAuth, ".")
}

func Authenticator(opts *config.Opts) http.HandlerFunc {
	switch opts.Auth {
	case anonAuth:
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "no authentication in anon mode", http.StatusBadRequest)
		})
	case sipAuth:
		return sip2.SipAuthenticator(opts)
	default:
		bailOnBadAuthModule(opts.Auth)
	}
	return nil
}

func RestrictedMiddleware(auth string, ch web.CertHandler) http.Handler {

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return jwtSigningSecret, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	switch auth {
	case anonAuth:
		return http.HandlerFunc(ch.CertResponder)
	case sipAuth:
		return jwtMiddleware.Handler(http.HandlerFunc(ch.CertResponder))
	default:
		bailOnBadAuthModule(auth)
	}
	return nil
}
