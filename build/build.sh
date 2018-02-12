#!/bin/sh

# build the kolumbus binary
cd ..
gox -osarch="linux/amd64" -output="examples/grpc/build/bin/kolumbus"
cd ./build
