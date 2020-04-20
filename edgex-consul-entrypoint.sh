#!/bin/sh
set -e

Database_Check="${EDGEX_DB:-redis}"
Secure_Check="${EDGEX_SECURE:-true}"

# Set default config
cp /edgex/config/00-consul.json /consul/config/00-consul.json

# Set config files - Redis DB
if [ "$Database_Check" = 'redis' ]; then
    echo "Installing redis health checks"
    cp /edgex/config/database/01-redis.json /consul/config/01-redis.json
fi

# Set config files - Mongo DB
if [ "$Database_Check" = 'mongo' ]; then
    echo "Installing mongo health checks"
    cp /edgex/config/database/01-mongo.json /consul/config/01-mongo.json
fi

# Set config files - Secure Setup
if [ "$Secure_Check" = 'true' ]; then
    echo "Installing security health checks"
    cp /edgex/config/secure/*.json /consul/config/
fi

# Copy health check scripts
cp -r /edgex/scripts/* /consul/scripts/

echo "Chaining to original entrypoint"
exec "docker-entrypoint.sh" "$@"