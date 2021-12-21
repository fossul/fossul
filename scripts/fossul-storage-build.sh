#!/bin/sh

if [[ -z "${STORAGE_PLUGIN_DIR}" ]]; then
    STORAGE_DIR="${HOME}/plugins/storage"
    export STORAGE_PLUGIN_DIR=$STORAGE_DIR

  if [[ ! -e "${STORAGE_PLUGIN_DIR}" ]]; then
      mkdir -p $STORAGE_PLUGIN_DIR
  fi
else
    export STORAGE_PLUGIN_DIR="."
fi

if [[ -z "${ARCHIVE_PLUGIN_DIR}" ]]; then
    ARCHIVE_DIR="${HOME}/plugins/archive"
    export ARCHIVE_PLUGIN_DIR=$ARCHIVE_DIR

  if [[ ! -e "${ARCHIVE_PLUGIN_DIR}" ]]; then
      mkdir -p $ARCHIVE_PLUGIN_DIR
  fi
else
    export ARCHIVE_PLUGIN_DIR="."
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
go build -buildmode=plugin -o $STORAGE_PLUGIN_DIR/sample-storage.so github.com/fossul/fossul/src/plugins/storage/native/sample-storage
if [ $? != 0 ]; then exit 1; fi
go install github.com/fossul/fossul/src/plugins/archive/basic/sample-archive
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $ARCHIVE_PLUGIN_DIR/sample-archive.so github.com/fossul/fossul/src/plugins/archive/native/sample-archive
if [ $? != 0 ]; then exit 1; fi
go install github.com/fossul/fossul/src/plugins/storage/basic/container-basic
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $STORAGE_PLUGIN_DIR/container-basic.so github.com/fossul/fossul/src/plugins/storage/native/container-basic
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $ARCHIVE_PLUGIN_DIR/aws.so github.com/fossul/fossul/src/plugins/archive/native/aws
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $STORAGE_PLUGIN_DIR/ocs-gluster.so github.com/fossul/fossul/src/plugins/storage/native/ocs-gluster
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $STORAGE_PLUGIN_DIR/csi-ceph.so github.com/fossul/fossul/src/plugins/storage/native/csi-ceph
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $STORAGE_PLUGIN_DIR/csi.so github.com/fossul/fossul/src/plugins/storage/native/csi
if [ $? != 0 ]; then exit 1; fi

echo "Building Storage Service"
go install github.com/fossul/fossul/src/engine/storage
if [ $? != 0 ]; then exit 1; fi

if [[ ! -z "${GOBIN}" ]]; then
	echo "Moving plugins to $PLUGIN_DIR"
	mv $GOBIN/sample-storage $STORAGE_PLUGIN_DIR
	if [ $? != 0 ]; then exit 1; fi
	mv $GOBIN/sample-archive $ARCHIVE_PLUGIN_DIR
	if [ $? != 0 ]; then exit 1; fi
	mv $GOBIN/container-basic $STORAGE_PLUGIN_DIR
	if [ $? != 0 ]; then exit 1; fi

	echo "Copying startup scripts"
	cp scripts/fossul-storage-startup.sh $GOBIN
	if [ $? != 0 ]; then exit 1; fi
fi

echo "Storage build completed successfully"

