#!/bin/sh -x


cd /consul/scripts || exit 1

# use JQ from path if it exists (i.e. on ubuntu), otherwise use from consul
# scripts dir with local LD_LIBRARY_PATH to load libonig from there too
JQ=$(command -v jq)
if [ -z "$JQ" ]; then
    LD_LIBRARY_PATH=.
    export LD_LIBRARY_PATH
    JQ=/consul/scripts/jq
fi
export JQ

if ! [ "true" = "$(curl -s -G edgex-core-consul:8500/v1/agent/checks | "$JQ" -r '.[] | select(.ServiceName == "'"$1"'") | .Status == "passing"')" ]; then
    exit 2
fi
