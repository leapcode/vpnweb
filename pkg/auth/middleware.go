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
const sip2Auth string = "sip"

func bailOnBadAuthModule(module string) {
	log.Fatal("Unknown auth module: '", module, "'. Should be one of: ", anonAuth, ", ", sipAuth, ".")
}

func checkForAuthSecret(opts *config.Opts) {
	if opts.AuthSecret == "" {
		log.Fatal("Need to provide a AuthSecret value for SIP Authentication")
	}
	if len(opts.AuthSecret) < 20 {
		log.Fatal("Please provider an AuthSecret longer than 20 chars")
	}
}

func AuthenticatorMiddleware(opts *config.Opts) http.HandlerFunc {
	switch opts.Auth {
	case anonAuth:
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "no authentication in anon mode", http.StatusBadRequest)
		})
	case sip2Auth:
		checkForAuthSecret(opts)
		return sip2.SipAuthenticator(opts)
	default:
		bailOnBadAuthModule(opts.Auth)
	}
	return nil
}

func RestrictedMiddleware(opts *config.Opts, ch web.CertHandler) http.Handler {

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(opts.AuthSecret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	switch opts.Auth {
	case anonAuth:
		return http.HandlerFunc(ch.CertResponder)
	case sip2Auth:
		checkForAuthSecret(opts)
		return jwtMiddleware.Handler(http.HandlerFunc(ch.CertResponder))
	default:
		bailOnBadAuthModule(opts.Auth)
	}
	return nil
}
