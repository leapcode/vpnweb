package main

import (
	"0xacab.org/leap/vpnweb/pkg/auth/sip2"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const authURI string = "http://%s:%s/3/auth"
const certURI string = "http://%s:%s/3/cert"

func formatCredentials(user, pass string) (string, error) {
	c := sip2.Credentials{user, pass}
	credJson, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(credJson), nil
}

func getToken(credJson, host, port string) string {
	resp, err := http.Post(fmt.Sprintf(authURI, host, port), "text/json", strings.NewReader(credJson))
	if err != nil {
		log.Fatal("Error on auth request: ", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 401 {
		log.Println("401 UNAUTHORIZED")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Cannot read response body")
	}
	return string(body)
}

func getCert(token, host, port string) string {
	req, err := http.NewRequest("POST", fmt.Sprintf(certURI, host, port), strings.NewReader(""))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("cannot read response body")
	}
	defer resp.Body.Close()
	if resp.StatusCode == 401 {
		log.Println("401 UNAUTHORIZED")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Cannot read response body")
	}
	return string(body)
}

func doAuthenticate(user, pass, host, port string) {
	credJson, err := formatCredentials(user, pass)
	if err != nil {
		log.Fatal("Cannot encode credentials: ", err)
	}
	token := getToken(credJson, host, port)
	log.Println("token:", token)
	cert := getCert(token, host, port)
	log.Println(cert)
}

func main() {
	var host, port, user, pass string
	flag.StringVar(&host, "host", "localhost", "Server to connect")
	flag.StringVar(&port, "port", "8000", "port to connect")
	flag.StringVar(&user, "user", "", "sip user to authenticate")
	flag.StringVar(&pass, "pass", "", "sip password to authenticate")
	flag.Parse()

	log.Println("connect to", host, port, "with credentials", user, ":", pass)
	doAuthenticate(user, pass, host, port)

}
