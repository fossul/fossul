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

echo "Building Plugins"
go install github.com/fossul/fossul/src/plugins/storage/basic/sample-storage
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/storage/sample-storage.so github.com/fossul/fossul/src/plugins/storage/native/sample-storage
if [ $? != 0 ]; then exit 1; fi
go install github.com/fossul/fossul/src/plugins/archive/basic/sample-archive
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/archive/sample-archive.so github.com/fossul/fossul/src/plugins/archive/native/sample-archive
if [ $? != 0 ]; then exit 1; fi
go install github.com/fossul/fossul/src/plugins/storage/basic/container-basic
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/storage/container-basic.so github.com/fossul/fossul/src/plugins/storage/native/container-basic
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/archive/aws.so github.com/fossul/fossul/src/plugins/archive/native/aws
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/storage/ocs-gluster.so github.com/fossul/fossul/src/plugins/storage/native/ocs-gluster
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/storage/csi-ceph.so github.com/fossul/fossul/src/plugins/storage/native/csi-ceph
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/storage/csi.so github.com/fossul/fossul/src/plugins/storage/native/csi
if [ $? != 0 ]; then exit 1; fi

echo "Building Storage Service"
go install github.com/fossul/fossul/src/engine/storage
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
cp $GOPATH/src/github.com/fossul/fossul/scripts/fossul-storage-startup.sh $GOBIN
if [ $? != 0 ]; then exit 1; fi

echo "Storage build completed successfully"

