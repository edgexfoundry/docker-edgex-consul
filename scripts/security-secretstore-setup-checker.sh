#!/bin/sh -ex

# vault will be reported as healthy when it is unsealed
# NOTE: comment out the vault healthy check below after we switch Vault to use 
# file backend as the check becomes superfluous
# the real check "SECRETSTORE_SETUP_DONE_FLAG" below already 
# waited on Vault successfully initialized at the first place
#/consul/scripts/consul-svc-healthy.sh vault

# If SECRETSTORE_SETUP_DONE_FLAG not defined, pass; else
# If SECRETSTORE_SETUP_DONE_FLAG file doesn't exist, error
test -z "${SECRETSTORE_SETUP_DONE_FLAG}" || test -f "${SECRETSTORE_SETUP_DONE_FLAG}"
