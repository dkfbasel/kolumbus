#!/bin/sh

# build the kolumbus binary
cd ..
gox -osarch="linux/amd64" -output="./build/kolumbus"
cp ./build/kolumbus ./examples/grpc-prod/build/bin/kolumbus
cd ./build

echo "-- build docker container (y/n):"
read buildContainer

if [ "$buildContainer" = "y" ] || [ "$buildContainer" = "" ]; then
  echo "specify container tag (dev):"
  read tag
  if [ "$tag" = "" ]; then
    tag="dev"
  fi
  docker build -t dkfbasel/kolumbus:$tag .

	echo "-- push docker container (y/n)"
	read pushContainer
	if [ "$pushContainer" = "y" ] || [ "$pushContainer" = "" ]; then
	  docker push dkfbasel/kolumbus:$tag
	fi

fi
