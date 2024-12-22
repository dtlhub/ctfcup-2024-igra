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
    AUTH_TOKEN="dtl:509bd7cc468edebcc2090088ea01b43116b4f3f68b993e4ec68652653c0281d7"
fi

AUTH_TOKEN="$AUTH_TOKEN" go run cmd/client/main.go -s "$ADDRESS" -l "$LEVEL"
