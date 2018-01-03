#!/bin/sh

/app/bin/service --port 8082 &
/usr/local/bin/envoy -c /app/envoy/service.json
