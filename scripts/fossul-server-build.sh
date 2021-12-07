#!/bin/sh

PLUGIN_DIR="${HOME}/plugins"

if [[ -z "${FOSSUL_BUILD_PLUGIN_DIR}" ]]; then
    export FOSSUL_BUILD_PLUGIN_DIR=$PLUGIN_DIR

  if [[ ! -e "${FOSSUL_BUILD_PLUGIN_DIR}" ]]; then
      mkdir $FOSSUL_BUILD_PLUGIN_DIR
  fi
fi  

echo "Installing Dependencies"
go mod tidy

echo "Running Unit Tests"
go test github.com/fossul/fossul/src/engine/util
if [ $? != 0 ]; then exit 1; fi
go test github.com/fossul/fossul/src/plugins/pluginUtil
if [ $? != 0 ]; then exit 1; fi

echo "Building Shared Libraries"
go build github.com/fossul/fossul/src/engine/util
if [ $? != 0 ]; then exit 1; fi
go build github.com/fossul/fossul/src/client
if [ $? != 0 ]; then exit 1; fi
go build github.com/fossul/fossul/src/client/k8s
if [ $? != 0 ]; then exit 1; fi
go build github.com/fossul/fossul/src/plugins/pluginUtil
if [ $? != 0 ]; then exit 1; fi

echo "Building Server Service"
go install github.com/fossul/fossul/src/engine/server
if [ $? != 0 ]; then exit 1; fi

echo "Server build completed successfully"

