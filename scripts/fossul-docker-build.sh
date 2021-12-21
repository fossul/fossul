#!/bin/bash

echo " Building docker images"

sudo podman build --squash -t fossul-server:latest -f src/engine/server .
if [ $? != 0 ]; then exit 1; fi
sudo podman build --squash -t fossul-app:latest -f src/engine/app .
if [ $? != 0 ]; then exit 1; fi
sudo podman build --squash -t fossul-storage:latest -f src/engine/storage .
if [ $? != 0 ]; then exit 1; fi

sudo podman login quay.io
if [ $? != 0 ]; then exit 1; fi

sudo push localhost/fossul-server:latest quay.io/ktenzer/fossul-server:latest
if [ $? != 0 ]; then exit 1; fi
sudo push localhost/fossul-server:latest quay.io/ktenzer/fossul-server:latest
if [ $? != 0 ]; then exit 1; fi
sudo push localhost/fossul-server:latest quay.io/ktenzer/fossul-server:latest
if [ $? != 0 ]; then exit 1; fi

echo "Building docker images completed successfully"
