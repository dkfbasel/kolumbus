#!/bin/sh

protoc \
  --go_out=plugins=grpc:./helloworld \
  helloworld.proto

protoc \
  --go_out=plugins=grpc:./helloworld2 \
  helloworld2.proto
