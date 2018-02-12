#!/bin/sh

# build the kolumbus binary
gox -osarch="linux/amd64" -output="examples/grpc/build/bin/kolumbus"
