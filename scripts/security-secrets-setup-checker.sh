#!/bin/sh -ex

# check that the root CA exists, the vault and kong certs exist
test -f /vault/config/pki/EdgeXFoundryCA/EdgeXFoundryCA.pem
test -f /vault/config/pki/EdgeXFoundryCA/edgex-vault.pem
test -f /vault/config/pki/EdgeXFoundryCA/edgex-kong.pem
