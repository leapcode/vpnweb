package main

import (
	"log"
	"net/http"

	"0xacab.org/leap/vpnweb/pkg/auth"
	"0xacab.org/leap/vpnweb/pkg/config"
	"0xacab.org/leap/vpnweb/pkg/web"
)

func main() {
	opts := config.NewOpts()
	ch := web.NewCertHandler(opts.CaCrt, opts.CaKey)

	/* protected routes */

	/* TODO https://0xacab.org/leap/vpnweb/issues/4
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

	if opts.tls == true {
		log.Fatal(http.ListenAndServeTLS(pstr, opts.TlsCrt, opts.TlsKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(pstr, nil))

	}
}
