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
	authenticator := auth.GetAuthenticator(opts, false)

	/* protected routes */

	/* TODO https://0xacab.org/leap/vpnweb/issues/4
	http.HandleFunc("/3/refresh-token", auth.RefreshAuthMiddleware(opts.Auth))
	*/

	http.HandleFunc("/3/auth", web.AuthMiddleware(authenticator.CheckCredentials, opts))
	http.Handle("/3/cert", web.RestrictedMiddleware(authenticator.NeedsCredentials, ch.CertResponder, opts))

	/* static files */

	web.HttpFileHandler("/3/configs.json", opts.ApiPath+"/3/configs.json")
	web.HttpFileHandler("/3/service.json", opts.ApiPath+"/3/service.json")
	web.HttpFileHandler("/3/config/eip-service.json", opts.ApiPath+"/3/eip-service.json")
	web.HttpFileHandler("/provider.json", opts.ApiPath+"provider.json")
	web.HttpFileHandler("/ca.crt", opts.ProviderCaPath)
	web.HttpFileHandler("/3/ca.crt", opts.ProviderCaPath)

	pstr := ":" + opts.Port
	log.Println("Listening in port", opts.Port)

	if opts.Tls == true {
		log.Fatal(http.ListenAndServeTLS(pstr, opts.TlsCrt, opts.TlsKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(pstr, nil))

	}
}
