package sip2

import (
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"

	"0xacab.org/leap/vpnweb/pkg/config"
)

const LibraryLocation string = "testlibrary"
const SipUser string = "leap"
const SipPasswd string = "Kohapassword1!"

// XXX duplicated, pass in opts
var jwtSigningSecret = []byte("thesingingkey")

type Credentials struct {
	User     string
	Password string
}

func SipAuthenticator(opts *config.Opts) http.HandlerFunc {
	log.Println("Initializing sip2 authenticator...")

	/* TODO -- should pass specific SIP options as a secondary struct */
	/* TODO -- catch connection errors */

	sip := NewClient("localhost", "6001", LibraryLocation)

	ok, err := sip.Connect()
	if err != nil {
		log.Fatal("cannot connect sip client")
	}
	ok = sip.Login(SipUser, SipPasswd)
	if !ok {
		log.Println("Error on SIP login")
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
		/* maybe no uid at all */
		claims["uid"] = "user"
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
		tokenString, _ := token.SignedString(jwtSigningSecret)
		w.Write([]byte(tokenString))
	})
	return authTokenHandler
}
