package main

import (
	"flag"
	"net/http"
)

const keySize = 2048
const expiryDays = 28

type certHandler struct {
	cainfo caInfo
}

func (ch *certHandler) certResponder(w http.ResponseWriter, r *http.Request) {
	ch.cainfo.CertWriter(w)
}

func main() {
	var caCrt = flag.String("caCrt", "", "path to the CA public key")
	var caKey = flag.String("caKey", "", "path to the CA private key")

	flag.Parse()

	ci := newCaInfo(*caCrt, *caKey)
	ch := certHandler{ci}

	http.HandleFunc("/cert", ch.certResponder)
	http.ListenAndServe(":8000", nil)
}
