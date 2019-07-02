package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// TODO get this from the config yaml?
const keySize = 2048
const expiryDays = 28

func errExit(errmsg string) {
	fmt.Printf("ERROR: %s\n", errmsg)
	os.Exit(1)
}

type certHandler struct {
	cainfo caInfo
}

func (ch *certHandler) certResponder(w http.ResponseWriter, r *http.Request) {
	ch.cainfo.CertWriter(w)
}

func doFilesSanityCheck(caCrt string, caKey string) {
	if _, err := os.Stat(caCrt); os.IsNotExist(err) {
		errExit("cannot find caCrt file")
	}
	if _, err := os.Stat(caKey); os.IsNotExist(err) {
		errExit("cannot find caKey file")
	}
}

func httpFileHandler(route string, path string) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	})
}

func main() {
	var caCrt = flag.String("caCrt", "", "path to the CA public key")
	var caKey = flag.String("caKey", "", "path to the CA private key")
	var port = flag.Int("port", 8000, "port where the server will listen")

	flag.Parse()

	if *caCrt == "" {
		errExit("missing caCrt parameter")
	}
	if *caKey == "" {
		errExit("missing caKey parameter")
	}

	doFilesSanityCheck(*caCrt, *caKey)

	ci := newCaInfo(*caCrt, *caKey)
	ch := certHandler{ci}

	// add routes here
	http.HandleFunc("/1/cert", ch.certResponder)
	httpFileHandler("/1/ca.crt", "./public/ca.crt")
	httpFileHandler("/1/configs.json", "./public/1/configs.json")
	httpFileHandler("/1/service.json", "./public/1/service.json")

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
