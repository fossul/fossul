#!/bin/sh

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
go build -buildmode=plugin -o sample-storage.so github.com/fossul/fossul/src/plugins/storage/native/sample-storage
if [ $? != 0 ]; then exit 1; fi
go install github.com/fossul/fossul/src/plugins/archive/basic/sample-archive
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o sample-archive.so github.com/fossul/fossul/src/plugins/archive/native/sample-archive
if [ $? != 0 ]; then exit 1; fi
go install github.com/fossul/fossul/src/plugins/storage/basic/container-basic
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o container-basic.so github.com/fossul/fossul/src/plugins/storage/native/container-basic
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o aws.so github.com/fossul/fossul/src/plugins/archive/native/aws
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o ocs-gluster.so github.com/fossul/fossul/src/plugins/storage/native/ocs-gluster
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o csi-ceph.so github.com/fossul/fossul/src/plugins/storage/native/csi-ceph
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o csi.so github.com/fossul/fossul/src/plugins/storage/native/csi
if [ $? != 0 ]; then exit 1; fi

echo "Building Storage Service"
go install github.com/fossul/fossul/src/engine/storage
if [ $? != 0 ]; then exit 1; fi

echo "Storage build completed successfully"

