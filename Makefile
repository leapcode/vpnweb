CONFIG=config/demo.yaml
EIP_TEMPLATE_3=scripts/templates/3/eip-service.json.jinja
EIP_TEMPLATE_4=scripts/templates/4/eip-service.json.jinja
EIP_SERVICE_3=deploy/public/3/eip-service.json
EIP_SERVICE_4=deploy/public/4/eip-service.json
PROVIDER_TEMPLATE=scripts/templates/provider.json.jinja
PROVIDER=deploy/public/provider.json

build:
	go build
demo-sip:
	. config/CONFIG && ./vpnweb -auth=sip2
demo-anon:
	. config/CONFIG && ./vpnweb -auth=anon
clean:
	rm -f public/1/*
	rm public/ca.crt
gen-shapeshifter:
	scripts/gen-shapeshifter-state.py deploy/shapeshifter-state
gen-provider:
	mkdir -p deploy/public/3
	mkdir -p deploy/public/4
	@python3 scripts/simplevpn.py --file=eip --config=$(CONFIG) --template=$(EIP_TEMPLATE_3) --obfs4_state deploy/shapeshifter-state > $(EIP_SERVICE_3) || echo "ERROR: see $(EIP_SERVICE_3) for output"
	@python3 scripts/simplevpn.py --file=eip --config=$(CONFIG) --template=$(EIP_TEMPLATE_4) --obfs4_state deploy/shapeshifter-state > $(EIP_SERVICE_4) || echo "ERROR: see $(EIP_SERVICE_4) for output"
	@python3 scripts/simplevpn.py --file=provider --config=$(CONFIG) --template=$(PROVIDER_TEMPLATE) > $(PROVIDER) || echo "ERROR: see $(PROVIDER) for output"
	rsync -ra deploy/public/ public/
populate:
	cp test/1/* public/1/
	cp test/files/ca.crt public/
