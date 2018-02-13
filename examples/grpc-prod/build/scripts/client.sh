#!/bin/sh

# Start the client service and connect to the grpc service via internal envoy
# proxy. Per default the envoy proxy should be configured to listen to internal
# communications on 127.0.0.1:8081
/app/bin/client &

# Start the envoy proxy with the given configuration (in v2 api) and log
# only level error and above
/usr/local/bin/envoy -c /app/envoy/client.json --v2-config-only -l error
