#!/bin/sh -e

if pg_isready -h kong-db -U kong -d kong; then
    exit 0
else
    exit 2
fi