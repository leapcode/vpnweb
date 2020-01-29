package sip2

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
	"time"

	"0xacab.org/leap/vpnweb/pkg/config"
)

const sipUserVar string = "VPNWEB_SIP_USER"
const sipPassVar string = "VPNWEB_SIP_PASS"
const sipPortVar string = "VPNWEB_SIP_PORT"
const sipHostVar string = "VPNWEB_SIP_HOST"
const sipLibrLocVar string = "VPNWEB_SIP_LIBR_LOCATION"
const sipTerminatorVar string = "VPNWEB_SIP_TERMINATOR"
const sipDefaultTerminator string = "\r\n"

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

func setupTerminatorFromEnv() {
	config.FallbackToEnv(&telnetTerminator, sipTerminatorVar, sipDefaultTerminator)
	if telnetTerminator == "\\r" {
		telnetTerminator = "\r"
	} else if telnetTerminator == "\\r\\n" {
		telnetTerminator = "\r\n"
	}
}

func SipAuthenticator(opts *config.Opts) http.HandlerFunc {

	log.Println("Initializing SIP2 authenticator")

	SipUser := getConfigFromEnv(sipUserVar)
	SipPass := getConfigFromEnv(sipPassVar)
	SipHost := getConfigFromEnv(sipHostVar)
	SipPort := getConfigFromEnv(sipPortVar)
	SipLibrLoc := getConfigFromEnv(sipLibrLocVar)

	setupTerminatorFromEnv()

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
