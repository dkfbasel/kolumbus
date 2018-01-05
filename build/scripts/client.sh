#!/bin/sh

# start the client service on port 8081 and connect to the grpc service
# via internal envoy proxy
/app/bin/client --port=8081 --grpc=127.0.0.1:9001 &

# start the envoy proxy with the given configuration (in v2 api) and log
# only level error and above
/usr/local/bin/envoy -c /app/envoy/client.json --v2-config-only -l error
