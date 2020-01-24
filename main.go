package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
)

// TODO get this from the config yaml?
const keySize = 2048
const expiryDays = 28

type certHandler struct {
	cainfo caInfo
}

func (ch *certHandler) certResponder(w http.ResponseWriter, r *http.Request) {
	ch.cainfo.CertWriter(w)
}

func doCaFilesSanityCheck(caCrt string, caKey string) {
	if _, err := os.Stat(caCrt); os.IsNotExist(err) {
		log.Fatal("cannot find caCrt file")
	}
	if _, err := os.Stat(caKey); os.IsNotExist(err) {
		log.Fatal("cannot find caKey file")
	}
}

func doTlsFilesSanityCheck(tlsCrt string, tlsKey string) {
	if _, err := os.Stat(tlsCrt); os.IsNotExist(err) {
		log.Fatal("cannot find tlsCrt file")
	}
	if _, err := os.Stat(tlsKey); os.IsNotExist(err) {
		log.Fatal("cannot find tlsKey file")
	}
}

func httpFileHandler(route string, path string) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	})
}

type Opts struct {
	caCrt  string
	caKey  string
	port   int
	notls  bool
	tlsCrt string
	tlsKey string
}

func initializeFlags(opts *Opts) {
	flag.StringVar(&opts.caCrt, "caCrt", "", "path to the CA public key")
	flag.StringVar(&opts.caKey, "caKey", "", "path to the CA private key")
	flag.IntVar(&opts.port, "port", 8000, "port where the server will listen")
	flag.BoolVar(&opts.notls, "notls", false, "disable TLS on the service")
	flag.StringVar(&opts.tlsCrt, "tls_crt", "", "path to the cert file for TLS")
	flag.StringVar(&opts.tlsKey, "tls_key", "", "path to the key file for TLS")
	flag.Parse()

	auth := os.Getenv("AUTH")
	log.Println("AUTH-->", auth)

}

func checkConfigurationOptions(opts *Opts) {

	if opts.caCrt == "" {
		log.Fatal("missing caCrt parameter")
	}
	if opts.caKey == "" {
		log.Fatal("missing caKey parameter")
	}

	if opts.notls == false {
		if opts.tlsCrt == "" {
			log.Fatal("missing tls_crt parameter. maybe use -notls?")
		}
		if opts.tlsKey == "" {
			log.Fatal("missing tls_key parameter. maybe use -notls?")
		}
	}

	doCaFilesSanityCheck(opts.caCrt, opts.caKey)
	if opts.notls == false {
		doTlsFilesSanityCheck(opts.tlsCrt, opts.tlsKey)
	}
}

func main() {
	opts := new(Opts)
	initializeFlags(opts)
	checkConfigurationOptions(opts)

	ci := newCaInfo(opts.caCrt, opts.caKey)
	ch := certHandler{ci}

	// add routes here
	http.HandleFunc("/3/cert", ch.certResponder)
	httpFileHandler("/3/configs.json", "./public/3/configs.json")
	httpFileHandler("/3/service.json", "./public/3/service.json")
	httpFileHandler("/3/config/eip-service.json", "./public/3/eip-service.json")
	httpFileHandler("/provider.json", "./public/provider.json")
	httpFileHandler("/ca.crt", "./public/ca.crt")
	httpFileHandler("/3/ca.crt", "./public/ca.crt")

	pstr := ":" + strconv.Itoa(opts.port)

	if opts.notls == true {
		log.Fatal(http.ListenAndServe(pstr, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(pstr, opts.tlsCrt, opts.tlsKey, nil))

	}
}
