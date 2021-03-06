#!/bin/sh

PLUGIN_DIR="/home/ktenzer/plugins"

if [[ -z "${FOSSUL_BUILD_PLUGIN_DIR}" ]]; then
    export FOSSUL_BUILD_PLUGIN_DIR=$PLUGIN_DIR
fi  

echo "Installing Dependencies"
$GOBIN/dep ensure

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

echo "Copying startup script"
cp $GOPATH/src/github.com/fossul/fossul/scripts/fossul-server-startup.sh $GOBIN
if [ $? != 0 ]; then exit 1; fi

echo "Server build completed successfully"

