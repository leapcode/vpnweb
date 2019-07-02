build:
	go build
demo:
	./vpnweb -caCrt test/files/ca.crt -caKey test/files/ca.key
clean:
	rm -f public/1/*
	rm public/ca.crt
populate:
	cp test/1/* public/1/
	cp test/files/ca.crt public/
