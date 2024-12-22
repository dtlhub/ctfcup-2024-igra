#!/usr/bin/env bash

set -e

if [ -z "$1" ]; then
	echo "Usage: $0 <local/remote>"
	exit 1
fi

LEVEL="level-3"

ADDRESS="localhost:8085"
AUTH_TOKEN="dtl:12345"
if [ "$1" == "remote" ]; then
    ADDRESS="172.16.30.195:13372"
    AUTH_TOKEN="dtl:42e11d41f32a24e58fb51db01fa2e8e5f0a22910f7e37b4de8c258e2ba7b5455"
fi

AUTH_TOKEN="$AUTH_TOKEN" go run cmd/client/main.go -s "$ADDRESS" -l "$LEVEL"
