#!/bin/sh

# build the service binary
cd ./src/service
gox -osarch="linux/amd64" -output="../../build/bin/service"
cd ../..

# build the client binary
cd ./src/client
gox -osarch="linux/amd64" -output="../../build/bin/client"
cd ../..
