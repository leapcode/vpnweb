package config

import (
	"flag"
	"log"
	"os"
	"reflect"
)

const DefaultAuthenticationModule = "anonymous"

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

func InitializeFlags(opts *Opts) {
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

func CheckConfigurationOptions(opts *Opts) {
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

	log.Println("Authentication module:", opts.Auth)

	// TODO -- check authentication module is valud, bail out otherwise
}
