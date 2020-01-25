package main

import (
	"log"
	"net/http"

	//"0xacab.org/leap/pkg/auth"
	"0xacab.org/leap/vpnweb/pkg/auth"
	"0xacab.org/leap/vpnweb/pkg/config"
	"0xacab.org/leap/vpnweb/pkg/web"
)

func main() {
	opts := new(config.Opts)
	config.InitializeFlags(opts)
	config.CheckConfigurationOptions(opts)

	ci := web.NewCaInfo(opts.CaCrt, opts.CaKey)
	ch := web.CertHandler{ci}

	/* protected routes */

	/* TODO ----
	http.HandleFunc("/3/refresh-token", auth.RefreshAuthMiddleware(opts.Auth))
	*/

	http.Handle("/3/cert", auth.RestrictedMiddleware(opts, ch))
	http.HandleFunc("/3/auth", auth.AuthenticatorMiddleware(opts))

	/* static files */

	/* TODO -- pass static file path in options */

	web.HttpFileHandler("/3/configs.json", "./public/3/configs.json")
	web.HttpFileHandler("/3/service.json", "./public/3/service.json")
	web.HttpFileHandler("/3/config/eip-service.json", "./public/3/eip-service.json")
	web.HttpFileHandler("/3/ca.crt", "./public/ca.crt")
	web.HttpFileHandler("/provider.json", "./public/provider.json")
	web.HttpFileHandler("/ca.crt", "./public/ca.crt")

	pstr := ":" + opts.Port
	log.Println("Listening in port", opts.Port)

	if opts.Notls == true {
		log.Fatal(http.ListenAndServe(pstr, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(pstr, opts.TlsCrt, opts.TlsKey, nil))

	}
}
