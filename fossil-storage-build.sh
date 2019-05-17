#!/bin/sh

PLUGIN_DIR="/home/ktenzer/plugins"

if [[ -z "${FOSSIL_BUILD_PLUGIN_DIR}" ]]; then
    export FOSSIL_BUILD_PLUGIN_DIR=$PLUGIN_DIR
fi  

echo "Installing Dependencies"
$GOBIN/dep ensure

echo "Running Unit Tests"
go test fossil/src/engine/util
if [ $? != 0 ]; then exit 1; fi
go test fossil/src/engine/plugins/pluginUtil
if [ $? != 0 ]; then exit 1; fi

echo "Building Shared Libraries"
go build fossil/src/engine/util
if [ $? != 0 ]; then exit 1; fi
go build fossil/src/engine/client
if [ $? != 0 ]; then exit 1; fi
go build fossil/src/engine/client/k8s
if [ $? != 0 ]; then exit 1; fi
go build fossil/src/engine/plugins/pluginUtil
if [ $? != 0 ]; then exit 1; fi

echo "Building Plugins"
go install fossil/src/engine/plugins/storage/basic/sample-storage
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/storage/sample-storage.so fossil/src/engine/plugins/storage/native/sample-storage
if [ $? != 0 ]; then exit 1; fi
go install fossil/src/engine/plugins/archive/basic/sample-archive
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/archive/sample-archive.so fossil/src/engine/plugins/archive/native/sample-archive
if [ $? != 0 ]; then exit 1; fi
go install fossil/src/engine/plugins/storage/basic/container-basic
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/storage/container-basic.so fossil/src/engine/plugins/storage/native/container-basic
if [ $? != 0 ]; then exit 1; fi

echo "Building Storage Service"
go install fossil/src/engine/storage
if [ $? != 0 ]; then exit 1; fi

echo "Moving plugins to $PLUGIN_DIR"
if [ ! -d "$PLUGIN_DIR/storage" ]; then mkdir $PLUGIN_DIR/storage; fi
mv $GOBIN/sample-storage $PLUGIN_DIR/storage
if [ $? != 0 ]; then exit 1; fi
if [ ! -d "$PLUGIN_DIR/archive" ]; then mkdir $PLUGIN_DIR/archive; fi
mv $GOBIN/sample-archive $PLUGIN_DIR/archive
if [ $? != 0 ]; then exit 1; fi
mv $GOBIN/container-basic $PLUGIN_DIR/storage
if [ $? != 0 ]; then exit 1; fi

echo "Copying startup script"
cp $GOPATH/src/fossil/fossil-storage-startup.sh $GOBIN
if [ $? != 0 ]; then exit 1; fi

echo "Storage build completed successfully"

