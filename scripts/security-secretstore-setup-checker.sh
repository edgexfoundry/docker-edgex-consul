#!/bin/sh -ex

# vault will be reported as healthy when it is unsealed
/consul/scripts/consul-svc-healthy.sh vault

# If SECRETSTORE_SETUP_DONE_FLAG not defined, pass; else
# If SECRETSTORE_SETUP_DONE_FLAG file doesn't exist, error
test -z "${SECRETSTORE_SETUP_DONE_FLAG}" || test -f "${SECRETSTORE_SETUP_DONE_FLAG}"
