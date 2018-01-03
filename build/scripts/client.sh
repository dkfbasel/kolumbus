#!/bin/sh

/app/bin/client --port=8081 --grpc=127.0.0.1:9001 &
/usr/local/bin/envoy -c /app/envoy/client.json --v2-config-only
