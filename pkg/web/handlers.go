package web

import (
	"net/http"
)

type CertHandler struct {
	Cainfo caInfo
}

func (ch *CertHandler) CertResponder(w http.ResponseWriter, r *http.Request) {
	ch.Cainfo.CertWriter(w)
}

func HttpFileHandler(route string, path string) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	})
}
