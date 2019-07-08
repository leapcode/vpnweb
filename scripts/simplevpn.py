#!/usr/bin/env python3
import argparse
import os

import yaml

from jinja2 import Template


class EIPConfig:
    def __init__(self):
        self.openvpn = dict()
        self.locations = dict()
        self.gateways = dict()
        self.obfs4_cert = ""


def parseConfig(provider_config):
    with open(provider_config) as conf:
        config = yaml.load(conf.read())
    eip = EIPConfig()
    eip.openvpn.update(yamlListToDict(config['openvpn']))

    for loc in config['locations']:
        eip.locations.update(yamlIdListToDict(loc))
    for gw in config['gateways']:
        eip.gateways.update(yamlIdListToDict(gw))
    return eip


def yamlListToDict(values):
    vals = {}
    for d in values:
        for k, v in d.items():
            vals[k] = v
    return vals


def yamlIdListToDict(data):
    _d = {}
    for identifier, values in data.items():
        _d[identifier] = yamlListToDict(values)
    return _d


def patchObfs4Cert(config, cert):
    for gw in config.gateways:
        for options in config.gateways[gw]['transports']:
            opts = {}
            transport, _, _ = options
            if transport == "obfs4":
                opts['cert'] = cert
                opts['iat-mode'] = 0
            options.append(opts)
    return config


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("provider_config")
    parser.add_argument("eip_template")
    parser.add_argument("--obfs4_state")
    args = parser.parse_args()

    config = parseConfig(os.path.abspath(args.provider_config))

    if args.obfs4_state:
        obfs4_cert = open(
            args.obfs4_state + '/obfs4_cert.txt').read().rstrip()
    else:
        obfs4_cert = None
    patchObfs4Cert(config, obfs4_cert)

    t = Template(open(args.eip_template).read())

    print(t.render(
        locations=config.locations,
        gateways=config.gateways,
        openvpn=config.openvpn))