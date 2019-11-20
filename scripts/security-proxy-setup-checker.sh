#!/bin/sh -ex

/consul/scripts/consul-svc-healthy.sh kong

# setup jq paths to check if kong has the required certificate uploaded to it
cd /consul/scripts || exit 1
LD_LIBRARY_PATH=.
export LD_LIBRARY_PATH
if [ "1" != "$(curl -s http://kong:8001/certificates | /consul/scripts/jq -r '.data | map(select(."snis" == ["edgex-kong"])) | length')" ]; then
    exit 2
fi
