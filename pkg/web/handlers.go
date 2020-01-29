package web

import (
	"net/http"
)

type CertHandler struct {
	Cainfo caInfo
}

func NewCertHandler(caCrt, caKey string) CertHandler {
	ci := newCaInfo(caCrt, caKey)
	ch := CertHandler{ci}
	return ch
}

func (ch *CertHandler) CertResponder(w http.ResponseWriter, r *http.Request) {
	ch.Cainfo.CertWriter(w)
}

func HttpFileHandler(route string, path string) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	})
}
