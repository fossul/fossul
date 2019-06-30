#!/bin/sh

PLUGIN_DIR="/home/ktenzer/plugins"

if [[ -z "${FOSSUL_BUILD_PLUGIN_DIR}" ]]; then
    export FOSSUL_BUILD_PLUGIN_DIR=$PLUGIN_DIR
fi  

echo "Installing Dependencies"
$GOBIN/dep ensure

echo "Running Unit Tests"
go test fossul/src/engine/util
if [ $? != 0 ]; then exit 1; fi
go test fossul/src/engine/plugins/pluginUtil
if [ $? != 0 ]; then exit 1; fi

echo "Building Shared Libraries"
go build fossul/src/engine/util
if [ $? != 0 ]; then exit 1; fi
go build fossul/src/engine/client
if [ $? != 0 ]; then exit 1; fi
go build fossul/src/engine/client/k8s
if [ $? != 0 ]; then exit 1; fi
go build fossul/src/engine/plugins/pluginUtil
if [ $? != 0 ]; then exit 1; fi

echo "Building Plugins"
go install fossul/src/engine/plugins/storage/basic/sample-storage
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/storage/sample-storage.so fossul/src/engine/plugins/storage/native/sample-storage
if [ $? != 0 ]; then exit 1; fi
go install fossul/src/engine/plugins/archive/basic/sample-archive
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/archive/sample-archive.so fossul/src/engine/plugins/archive/native/sample-archive
if [ $? != 0 ]; then exit 1; fi
go install fossul/src/engine/plugins/storage/basic/container-basic
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/storage/container-basic.so fossul/src/engine/plugins/storage/native/container-basic
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/archive/aws.so fossul/src/engine/plugins/archive/native/aws
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/storage/ocs-gluster.so fossul/src/engine/plugins/storage/native/ocs-gluster
if [ $? != 0 ]; then exit 1; fi

echo "Building Storage Service"
go install fossul/src/engine/storage
if [ $? != 0 ]; then exit 1; fi

echo "Moving plugins to $FOSSUL_BUILD_PLUGIN_DIR"
if [ ! -d "$FOSSUL_BUILD_PLUGIN_DIR/storage" ]; then mkdir $FOSSUL_BUILD_PLUGIN_DIR/storage; fi
mv $GOBIN/sample-storage $FOSSUL_BUILD_PLUGIN_DIR/storage
if [ $? != 0 ]; then exit 1; fi
if [ ! -d "$FOSSUL_BUILD_PLUGIN_DIR/archive" ]; then mkdir $FOSSUL_BUILD_PLUGIN_DIR/archive; fi
mv $GOBIN/sample-archive $FOSSUL_BUILD_PLUGIN_DIR/archive
if [ $? != 0 ]; then exit 1; fi
mv $GOBIN/container-basic $FOSSUL_BUILD_PLUGIN_DIR/storage
if [ $? != 0 ]; then exit 1; fi

echo "Copying startup script"
cp $GOPATH/src/fossul/scripts/fossul-storage-startup.sh $GOBIN
if [ $? != 0 ]; then exit 1; fi

echo "Storage build completed successfully"

