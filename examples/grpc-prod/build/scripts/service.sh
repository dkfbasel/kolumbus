#!/bin/sh

# Start the service and connect and internal envoy proxy. Per default the
# envoy proxy should be configured to listen to incoming requests on
# 0.0.0.0:80 and forward it locally to 127.0.0.1:8080

/app/bin/service &
/usr/local/bin/envoy -c /app/envoy/service.json --v2-config-only -l error
