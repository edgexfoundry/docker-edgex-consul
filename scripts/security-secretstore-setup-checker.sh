#!/bin/sh -ex

# for now check that the root token exists
# TODO: check that individual tokens exist when we don't use the root token 
# anymore
test -f /vault/config/assets/resp-init.json

/consul/scripts/consul-svc-healthy.sh security-secrets-setup

# vault will be reported as healthy when it is unsealed
/consul/scripts/consul-svc-healthy.sh vault
