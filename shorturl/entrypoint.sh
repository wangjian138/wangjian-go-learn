#!/bin/sh
cd /usr/local/shorturl/
go mod tidy

go run /usr/local/shorturl/rpc/transform/transform.go -f /usr/local/shorturl/rpc/transform/etc/transform.yaml

go run /usr/local/shorturl/api/shorturl.go  --consulAddr="106.13.191.41:8500"


exec "$@"