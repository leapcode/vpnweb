CERTTOOL=certtool
CAKEY=ca-key.pem
CATEMPL=ca.cfg
SRVTEMPL=server-cert.cfg
DEPLOY=certs

usage:
	@echo "Use 'make ca' for creating a new disposable ca, or 'make cert' for creating a new cert"

clean:
	rm -f tls-key.pem tls-cert.pem ca-key.pem ca.crt request.pem

ca:
	$(CERTTOOL) --generate-privkey --outfile $(CAKEY)
	$(CERTTOOL) --generate-self-signed --load-privkey ca-key.pem --outfile ca.crt --template $(CATEMPL)

cert:
	$(CERTTOOL) --generate-privkey --outfile tls-key.pem --template $(SRVTEMPL)
	$(CERTTOOL) --generate-request --load-privkey tls-key.pem --outfile request.pem --template $(SRVTEMPL)
	$(CERTTOOL) --generate-certificate --load-privkey tls-key.pem \
		    --template $(SRVTEMPL) --outfile tls-cert.pem \
		    --load-ca-certificate ca.crt --load-ca-privkey $(CAKEY)
deploy:
	@rm -rf $(DEPLOY)
	@mkdir $(DEPLOY)
	@cp ca.crt tls-key.pem tls-cert.pem $(DEPLOY)
	@echo "done"

