#!/bin/sh

# build the kolumbus binary
cd ..
gox -osarch="linux/amd64" -output="./build/kolumbus"
cp ./build/kolumbus ./examples/grpc/build/bin/kolumbus
cd ./build
