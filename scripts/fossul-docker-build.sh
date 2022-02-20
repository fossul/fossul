#!/bin/bash
# Needs to be run from fossul root directory

RELEASE="v0.5.0"
SERVER_IMAGE="fossul-server"
APP_IMAGE="fossul-app"
STORAGE_IMAGE="fossul-storage"
CLI_IMAGE="fossul-cli"
REPO="fossul"

if [[ -z "${BUILD_ARGS}" ]]; then
	export BUILD_ARGS="all"
fi

echo " Building docker images"

case $BUILD_ARGS in
  all)
    podman build -t quay.io/$REPO/$SERVER_IMAGE:$RELEASE -f src/engine/server .
    if [ $? != 0 ]; then exit 1; fi
    podman build -t quay.io/$REPO/$APP_IMAGE:$RELEASE -f src/engine/app .
    if [ $? != 0 ]; then exit 1; fi
    podman build -t quay.io/$REPO/$STORAGE_IMAGE:$RELEASE -f src/engine/storage .
    if [ $? != 0 ]; then exit 1; fi
    podman build -t quay.io/$REPO/$CLI_IMAGE:$RELEASE -f src/cli .
    if [ $? != 0 ]; then exit 1; fi

    podman tag quay.io/$REPO/$SERVER_IMAGE:$RELEASE quay.io/$REPO/$SERVER_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi
    podman tag quay.io/$REPO/$APP_IMAGE:$RELEASE quay.io/$REPO/$APP_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi
    podman tag quay.io/$REPO/$STORAGE_IMAGE:$RELEASE quay.io/$REPO/$STORAGE_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi
    podman tag quay.io/$REPO/$CLI_IMAGE:$RELEASE quay.io/$REPO/$CLI_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi

    podman login quay.io
    if [ $? != 0 ]; then exit 1; fi

    podman push quay.io/$REPO/$SERVER_IMAGE:$RELEASE
    if [ $? != 0 ]; then exit 1; fi
    podman push quay.io/$REPO/$APP_IMAGE:$RELEASE
    if [ $? != 0 ]; then exit 1; fi
    podman push quay.io/$REPO/$STORAGE_IMAGE:$RELEASE
    if [ $? != 0 ]; then exit 1; fi
    podman push quay.io/$REPO/$CLI_IMAGE:$RELEASE
    if [ $? != 0 ]; then exit 1; fi

    podman push quay.io/$REPO/$SERVER_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi
    podman push quay.io/$REPO/$APP_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi
    podman push quay.io/$REPO/$STORAGE_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi
    podman push quay.io/$REPO/$CLI_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi

    echo "Building docker images completed successfully"
    ;;
  server)
    podman build -t quay.io/$REPO/$SERVER_IMAGE:$RELEASE -f src/engine/server .
    if [ $? != 0 ]; then exit 1; fi

    podman tag quay.io/$REPO/$SERVER_IMAGE:$RELEASE quay.io/$REPO/$SERVER_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi

    podman login quay.io
    if [ $? != 0 ]; then exit 1; fi

    podman push quay.io/$REPO/$SERVER_IMAGE:$RELEASE
    if [ $? != 0 ]; then exit 1; fi

    podman push quay.io/$REPO/$SERVER_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi

    echo "Building server docker image completed successfully"
    ;;
  app)
    podman build -t quay.io/$REPO/$APP_IMAGE:$RELEASE -f src/engine/app .
    if [ $? != 0 ]; then exit 1; fi

    podman tag quay.io/$REPO/$APP_IMAGE:$RELEASE quay.io/$REPO/$APP_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi

    podman login quay.io
    if [ $? != 0 ]; then exit 1; fi

    podman push quay.io/$REPO/$APP_IMAGE:$RELEASE
    if [ $? != 0 ]; then exit 1; fi

    podman push quay.io/$REPO/$APP_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi

    echo "Building app docker image completed successfully"
    ;;
  storage)
    podman build -t quay.io/$REPO/$STORAGE_IMAGE:$RELEASE -f src/engine/storage .
    if [ $? != 0 ]; then exit 1; fi

    podman tag quay.io/$REPO/$STORAGE_IMAGE:$RELEASE quay.io/$REPO/$STORAGE_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi

    podman login quay.io
    if [ $? != 0 ]; then exit 1; fi

    podman push quay.io/$REPO/$STORAGE_IMAGE:$RELEASE
    if [ $? != 0 ]; then exit 1; fi

    podman push quay.io/$REPO/$STORAGE_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi

    echo "Building storage docker image completed successfully"
    ;;
  cli)
    podman build -t quay.io/$REPO/$CLI_IMAGE:$RELEASE -f src/cli .
    if [ $? != 0 ]; then exit 1; fi

    podman tag quay.io/$REPO/$CLI_IMAGE:$RELEASE quay.io/$REPO/$CLI_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi

    podman login quay.io
    if [ $? != 0 ]; then exit 1; fi

    podman push quay.io/$REPO/$CLI_IMAGE:$RELEASE
    if [ $? != 0 ]; then exit 1; fi

    podman push quay.io/$REPO/$CLI_IMAGE:latest
    if [ $? != 0 ]; then exit 1; fi

    echo "Building cli docker image completed successfully"
    ;;
  *)
    echo "Invalid build argument!"
    exit 1
    ;;
esac
