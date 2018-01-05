#!/bin/sh
cd ./src/service
gox -osarch="linux/amd64" -output="../../build/bin/service"

cd ../client
gox -osarch="linux/amd64" -output="../../build/bin/client"

cd ../dns
gox -osarch="linux/amd64" -output="../../build/bin/dns"

cd ../..
