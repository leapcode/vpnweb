package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"reflect"
)

const keySize = 2048
const expiryDays = 28
const DefaultAuthenticationModule = "anonymous"

type certHandler struct {
	cainfo caInfo
}

func (ch *certHandler) certResponder(w http.ResponseWriter, r *http.Request) {
	ch.cainfo.CertWriter(w)
}

type Opts struct {
	Notls  bool
	CaCrt  string
	CaKey  string
	TlsCrt string
	TlsKey string
	Port   string
	Auth   string
}

func (o *Opts) fallbackToEnv(field string, envVar string, defaultVal string) {
	r := reflect.ValueOf(o)
	f := reflect.Indirect(r).FieldByName(field)

	if f.String() == "" {
		val, exists := os.LookupEnv(envVar)
		if exists && val != "" {
			f.SetString(val)
		} else {
			f.SetString(defaultVal)
		}
	}
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

func initializeFlags(opts *Opts) {

	flag.BoolVar(&opts.Notls, "notls", false, "disable TLS on the service")
	flag.StringVar(&opts.CaCrt, "caCrt", "", "path to the CA public key")
	flag.StringVar(&opts.CaKey, "caKey", "", "path to the CA private key")
	flag.StringVar(&opts.TlsCrt, "tls_crt", "", "path to the cert file for TLS")
	flag.StringVar(&opts.TlsKey, "tls_key", "", "path to the key file for TLS")
	flag.StringVar(&opts.Port, "port", "", "port where the server will listen (default: 8000)")
	flag.StringVar(&opts.Auth, "auth", "", "authentication module (anonymous, sip)")
	flag.Parse()

	opts.fallbackToEnv("CaCrt", "VPNWEB_CACRT", "")
	opts.fallbackToEnv("CaKey", "VPNWEB_CAKEY", "")
	opts.fallbackToEnv("TlsCrt", "VPNWEB_TLSCRT", "")
	opts.fallbackToEnv("TlsKey", "VPNWEB_TLSKEY", "")
	opts.fallbackToEnv("Port", "VPNWEB_PORT", "8000")
	opts.fallbackToEnv("Auth", "VPNWEB_AUTH", DefaultAuthenticationModule)
}

func checkConfigurationOptions(opts *Opts) {

	if opts.CaCrt == "" {
		log.Fatal("missing caCrt parameter")
	}
	if opts.CaKey == "" {
		log.Fatal("missing caKey parameter")
	}

	if opts.Notls == false {
		if opts.TlsCrt == "" {
			log.Fatal("missing tls_crt parameter. maybe use -notls?")
		}
		if opts.TlsKey == "" {
			log.Fatal("missing tls_key parameter. maybe use -notls?")
		}
	}

	doCaFilesSanityCheck(opts.CaCrt, opts.CaKey)
	if opts.Notls == false {
		doTlsFilesSanityCheck(opts.TlsCrt, opts.TlsKey)
	}

	log.Println("authentication module:", opts.Auth)

	// TODO -- check authentication module is valud, bail out otherwise
}

func main() {
	opts := new(Opts)
	initializeFlags(opts)
	checkConfigurationOptions(opts)

	ci := newCaInfo(opts.CaCrt, opts.CaKey)
	ch := certHandler{ci}

	// add routes here
	http.HandleFunc("/3/cert", ch.certResponder)
	httpFileHandler("/3/configs.json", "./public/3/configs.json")
	httpFileHandler("/3/service.json", "./public/3/service.json")
	httpFileHandler("/3/config/eip-service.json", "./public/3/eip-service.json")
	httpFileHandler("/provider.json", "./public/provider.json")
	httpFileHandler("/ca.crt", "./public/ca.crt")
	httpFileHandler("/3/ca.crt", "./public/ca.crt")

	pstr := ":" + opts.Port
	log.Println("serving vpnweb in port", opts.Port)

	if opts.Notls == true {
		log.Fatal(http.ListenAndServe(pstr, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(pstr, opts.TlsCrt, opts.TlsKey, nil))

	}
}
