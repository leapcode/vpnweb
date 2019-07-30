build:
	go build
demo:
	./vpnweb -caCrt test/files/ca.crt -caKey test/files/ca.key
clean:
	rm -f public/1/*
	rm public/ca.crt
gen-shapeshifter:
	scripts/gen-shapeshifter-state.py deploy/shapeshifter-state
gen-provider:
	mkdir -p deploy/public/3
	python3 scripts/simplevpn.py config/demo.yaml scripts/templates/eip-service.json.jinja --obfs4_state deploy/shapeshifter-state > deploy/public/3/eip-service.json
populate:
	cp test/1/* public/1/
	cp test/files/ca.crt public/
