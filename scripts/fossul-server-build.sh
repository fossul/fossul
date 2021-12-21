#!/bin/bash

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

if [[ ! -z "${GOBIN}" ]]; then
	echo "Copying startup scripts"
	cp scripts/fossul-server-startup.sh $GOBIN
	if [ $? != 0 ]; then exit 1; fi
fi

echo "Server build completed successfully"

