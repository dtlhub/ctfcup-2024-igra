#!/usr/bin/env bash

set -e

if [ -z "$1" ]; then
	echo "Usage: $0 <local/remote>"
	exit 1
fi

LEVEL="level-2"

ADDRESS="localhost:8085"
AUTH_TOKEN="dtl:12345"
if [ "$1" == "remote" ]; then
    ADDRESS="172.16.30.195:13372"
    AUTH_TOKEN="dtl:794412228410f7209a7700b9eeb9f3203fb1c9e4bd87f3eba2a244489aa32395"
fi

AUTH_TOKEN="$AUTH_TOKEN" go run cmd/client/main.go -s "$ADDRESS" -l "$LEVEL"
