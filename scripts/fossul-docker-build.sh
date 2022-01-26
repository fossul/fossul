#!/bin/bash
# Needs to be run from fossul root directory

SERVER_IMAGE="fossul-server:latest"
APP_IMAGE="fossul-app:latest"
STORAGE_IMAGE="fossul-storage:latest"
CLI_IMAGE="fossul-cli:latest"
REPO="fossul"

if [[ -z "${BUILD_ARGS}" ]]; then
	export BUILD_ARGS="all"
fi

echo " Building docker images"

case $BUILD_ARGS in
  all)
    sudo podman build -t quay.io/$REPO/$SERVER_IMAGE -f src/engine/server .
    if [ $? != 0 ]; then exit 1; fi
    sudo podman build -t quay.io/$REPO/$APP_IMAGE -f src/engine/app .
    if [ $? != 0 ]; then exit 1; fi
    sudo podman build -t quay.io/$REPO/$STORAGE_IMAGE -f src/engine/storage .
    if [ $? != 0 ]; then exit 1; fi
    sudo podman build -t quay.io/$REPO/$CLI_IMAGE -f src/cli .
    if [ $? != 0 ]; then exit 1; fi

    sudo podman login quay.io
    if [ $? != 0 ]; then exit 1; fi

    sudo podman push quay.io/$REPO/$SERVER_IMAGE
    if [ $? != 0 ]; then exit 1; fi
    sudo podman push quay.io/$REPO/$APP_IMAGE
    if [ $? != 0 ]; then exit 1; fi
    sudo podman push quay.io/$REPO/$STORAGE_IMAGE
    if [ $? != 0 ]; then exit 1; fi
    sudo podman push quay.io/$REPO/$CLI_IMAGE
    if [ $? != 0 ]; then exit 1; fi

    echo "Building docker images completed successfully"
    ;;
  server)
    sudo podman build -t quay.io/$REPO/$SERVER_IMAGE -f src/engine/server .
    if [ $? != 0 ]; then exit 1; fi

    sudo podman login quay.io
    if [ $? != 0 ]; then exit 1; fi

    sudo podman push quay.io/$REPO/$SERVER_IMAGE
    if [ $? != 0 ]; then exit 1; fi

    echo "Building server docker image completed successfully"
    ;;
  app)
    sudo podman build -t quay.io/$REPO/$APP_IMAGE -f src/engine/app .
    if [ $? != 0 ]; then exit 1; fi

    sudo podman login quay.io
    if [ $? != 0 ]; then exit 1; fi

    sudo podman push quay.io/$REPO/$APP_IMAGE
    if [ $? != 0 ]; then exit 1; fi

    echo "Building app docker image completed successfully"
    ;;
  storage)
    sudo podman build -t quay.io/$REPO/$STORAGE_IMAGE -f src/engine/storage .
    if [ $? != 0 ]; then exit 1; fi

    sudo podman login quay.io
    if [ $? != 0 ]; then exit 1; fi

    sudo podman push quay.io/$REPO/$STORAGE_IMAGE
    if [ $? != 0 ]; then exit 1; fi

    echo "Building storage docker image completed successfully"
    ;;
  cli)
    sudo podman build -t quay.io/$REPO/$CLI_IMAGE -f src/cli .
    if [ $? != 0 ]; then exit 1; fi

    sudo podman login quay.io
    if [ $? != 0 ]; then exit 1; fi

    sudo podman push quay.io/$REPO/$CLI_IMAGE
    if [ $? != 0 ]; then exit 1; fi

    echo "Building cli docker image completed successfully"
    ;;
  *)
    echo "Invalid build argument!"
    exit 1
    ;;
esac
