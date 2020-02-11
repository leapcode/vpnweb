package main

import (
	"log"
	"net/http"

	"0xacab.org/leap/vpnweb/pkg/auth"
	"0xacab.org/leap/vpnweb/pkg/config"
	"0xacab.org/leap/vpnweb/pkg/web"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	opts := config.NewOpts()
	ch := web.NewCertHandler(opts.CaCrt, opts.CaKey)
	authenticator := auth.GetAuthenticator(opts, false)

	srv := http.NewServeMux()

	/* protected routes */

	/* TODO https://0xacab.org/leap/vpnweb/issues/4
	http.HandleFunc("/3/refresh-token", auth.RefreshAuthMiddleware(opts.Auth))
	*/
	srv.HandleFunc("/3/auth", web.AuthMiddleware(authenticator.CheckCredentials, opts))
	srv.Handle("/3/cert", web.RestrictedMiddleware(authenticator.NeedsCredentials, ch.CertResponder, opts))

	/* static files */

	web.HttpFileHandler(srv, "/3/configs.json", opts.ApiPath+"/3/configs.json")
	web.HttpFileHandler(srv, "/3/service.json", opts.ApiPath+"/3/service.json")
	web.HttpFileHandler(srv, "/3/config/eip-service.json", opts.ApiPath+"/3/eip-service.json")
	web.HttpFileHandler(srv, "/provider.json", opts.ApiPath+"provider.json")
	web.HttpFileHandler(srv, "/ca.crt", opts.ProviderCaPath)
	web.HttpFileHandler(srv, "/3/ca.crt", opts.ProviderCaPath)

	mtr := http.NewServeMux()
	mtr.Handle("/metrics", promhttp.Handler())

	/* prometheus metrics */
	go func() {
		pstr := ":" + opts.MetricsPort
		log.Println("/metrics endpoint in port", opts.MetricsPort)
		log.Fatal(http.ListenAndServe(pstr, mtr))
	}()

	/* api server */
	pstr := ":" + opts.Port
	log.Println("API listening in port", opts.Port)
	if opts.Tls == true {
		log.Fatal(http.ListenAndServeTLS(pstr, opts.TlsCrt, opts.TlsKey, srv))
	} else {
		log.Fatal(http.ListenAndServe(pstr, srv))
	}
}
