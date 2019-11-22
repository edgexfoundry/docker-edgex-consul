#!/bin/sh -ex

# check that the security-secrets-setup complete
# the new way of checking just the sentinel files which will have all TLS assets
test -f /tmp/edgex/secrets/ca/.security-secrets-setup.complete
