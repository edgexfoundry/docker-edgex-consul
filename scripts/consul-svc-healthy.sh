#!/bin/sh -x


cd /consul/scripts || exit 1

LD_LIBRARY_PATH=.
export LD_LIBRARY_PATH
if ! [ "true" = "$(curl -s -G edgex-core-consul:8500/v1/agent/checks | /consul/scripts/jq -r '.[] | select(.ServiceName == "'"$1"'") | .Status == "passing"')" ]; then
    exit 2
fi
