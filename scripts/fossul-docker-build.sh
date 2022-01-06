#!/bin/bash

SERVER_IMAGE="fossul-server:latest"
APP_IMAGE="fossul-app:latest"
STORAGE_IMAGE="fossul-storage:latest"
CLI_IMAGE="fossul-cli:latest"

echo " Building docker images"

sudo podman build -t quay.io/ktenzer/$SERVER_IMAGE -f src/engine/server .
if [ $? != 0 ]; then exit 1; fi
sudo podman build -t quay.io/ktenzer/$APP_IMAGE -f src/engine/app .
if [ $? != 0 ]; then exit 1; fi
sudo podman build -t quay.io/ktenzer/$STORAGE_IMAGE -f src/engine/storage .
if [ $? != 0 ]; then exit 1; fi
sudo podman build -t quay.io/ktenzer/$CLI_IMAGE -f src/cli .
if [ $? != 0 ]; then exit 1; fi

sudo podman login quay.io
if [ $? != 0 ]; then exit 1; fi

sudo podman push quay.io/ktenzer/$SERVER_IMAGE
if [ $? != 0 ]; then exit 1; fi
sudo podman push quay.io/ktenzer/$APP_IMAGE
if [ $? != 0 ]; then exit 1; fi
sudo podman push quay.io/ktenzer/$STORAGE_IMAGE
if [ $? != 0 ]; then exit 1; fi
sudo podman push quay.io/ktenzer/$CLI_IMAGE
if [ $? != 0 ]; then exit 1; fi

echo "Building docker images completed successfully"
