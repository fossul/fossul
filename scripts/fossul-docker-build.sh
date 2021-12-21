#!/bin/bash

SERVER_IMAGE="fossul-server:latest"
APP_IMAGE="fossul-app:latest"
STORAGE_IMAGE="fossul-storage:latest"

echo " Building docker images"

#sudo podman build --squash -t $SERVER_IMAGE -f src/engine/server .
#if [ $? != 0 ]; then exit 1; fi
sudo podman build --squash -t $APP_IMAGE -f src/engine/app .
if [ $? != 0 ]; then exit 1; fi
sudo podman build --squash -t $STORAGE_IMAGE -f src/engine/storage .
if [ $? != 0 ]; then exit 1; fi

sudo podman login quay.io
if [ $? != 0 ]; then exit 1; fi

sudo podman push localhost/$SERVER_IMAGE quay.io/ktenzer/$SERVER_IMAGE
if [ $? != 0 ]; then exit 1; fi
sudo podman push localhost/$APP_IMAGE quay.io/ktenzer/$APP_IMAGE
if [ $? != 0 ]; then exit 1; fi
sudo podman push localhost/$STORAGE_IMAGE quay.io/ktenzer/$STORAGE_IMAGE
if [ $? != 0 ]; then exit 1; fi

echo "Building docker images completed successfully"
