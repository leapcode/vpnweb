package sip2

import (
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"time"

	"0xacab.org/leap/vpnweb/pkg/config"
)

const SipUserVar string = "VPNWEB_SIP_USER"
const SipPassVar string = "VPNWEB_SIP_PASS"
const SipPortVar string = "VPNWEB_SIP_PORT"
const SipHostVar string = "VPNWEB_SIP_HOST"
const SipLibrLocVar string = "VPNWEB_SIP_LIBR_LOCATION"

type Credentials struct {
	User     string
	Password string
}

func getConfigFromEnv(envVar string) string {
	val, exists := os.LookupEnv(envVar)
	if !exists {
		log.Fatal("Need to set required env var:", envVar)
	}
	return val
}

func SipAuthenticator(opts *config.Opts) http.HandlerFunc {
	/* TODO -- catch connection errors */

	log.Println("Initializing sip2 authenticator")

	SipUser := getConfigFromEnv(SipUserVar)
	SipPass := getConfigFromEnv(SipPassVar)
	SipHost := getConfigFromEnv(SipHostVar)
	SipPort := getConfigFromEnv(SipPortVar)
	SipLibrLoc := getConfigFromEnv(SipLibrLocVar)

	sip := NewClient(SipHost, SipPort, SipLibrLoc)

	ok, err := sip.Connect()
	if err != nil {
		log.Fatal("Cannot connect sip client")
	}
	ok = sip.Login(SipUser, SipPass)
	if !ok {
		log.Fatal("Error on SIP login")
	} else {
		log.Println("SIP login ok")
	}

	var authTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var c Credentials

		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			log.Println("Auth request did not send valid json")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if c.User == "" || c.Password == "" {
			log.Println("Auth request did not include user or password")
			http.Error(w, "missing user and/or password", http.StatusBadRequest)
			return
		}

		valid := sip.CheckCredentials(c.User, c.Password)
		if !valid {
			log.Println("Wrong auth for user", c.User)
			http.Error(w, "wrong user and/or password", http.StatusUnauthorized)
			return
		}

		log.Println("Valid auth for user", c.User)
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		tokenString, _ := token.SignedString([]byte(opts.AuthSecret))
		w.Write([]byte(tokenString))
	})
	return authTokenHandler
}
