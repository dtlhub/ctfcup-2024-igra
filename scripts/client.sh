#!/usr/bin/env bash

set -e

if [ -z "$1" ]; then
	echo "Usage: $0 <local/remote>"
	exit 1
fi

LEVEL="test"

ADDRESS="localhost:8085"
if [ "$1" == "local" ]; then
    ADDRESS="localhost:8085"
else
    echo "TODO: set remote address"
    exit 1
fi

AUTH_TOKEN=dtl:12345 go run cmd/client/main.go -s "$ADDRESS" -l "$LEVEL"
