build:
	go build
demo:
	./vpnweb -caCrt test/files/ca.crt -caKey test/files/ca.key
clean:
	rm -f public/1/*
populate:
	cp test/1/* public/1/
