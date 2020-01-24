build:
	go build
demo:
	. config/CONFIG && ./vpnweb -caCrt test/files/ca.crt -caKey test/files/ca.key -notls
clean:
	rm -f public/1/*
	rm public/ca.crt
gen-shapeshifter:
	scripts/gen-shapeshifter-state.py deploy/shapeshifter-state
gen-provider:
	mkdir -p deploy/public/3
	python3 scripts/simplevpn.py --file=eip --config=config/demo.yaml --template=scripts/templates/eip-service.json.jinja --obfs4_state deploy/shapeshifter-state > deploy/public/3/eip-service.json
	python3 scripts/simplevpn.py --file=provider --config=config/demo.yaml --template=scripts/templates/provider.json.jinja > deploy/public/provider.json
populate:
	cp test/1/* public/1/
	cp test/files/ca.crt public/
