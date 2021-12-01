#!/bin/sh

PLUGIN_DIR="/home/ktenzer/plugins"

if [[ -z "${FOSSUL_BUILD_PLUGIN_DIR}" ]]; then
    export FOSSUL_BUILD_PLUGIN_DIR=$PLUGIN_DIR
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
go install github.com/fossul/fossul/src/plugins/app/basic/sample-app
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/app/sample-app.so github.com/fossul/fossul/src/plugins/app/native/sample-app
if [ $? != 0 ]; then exit 1; fi
go install github.com/fossul/fossul/src/plugins/storage/basic/sample-storage
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/storage/sample-storage.so github.com/fossul/fossul/src/plugins/storage/native/sample-storage
if [ $? != 0 ]; then exit 1; fi
go install github.com/fossul/fossul/src/plugins/archive/basic/sample-archive
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/archive/sample-archive.so github.com/fossul/fossul/src/plugins/archive/native/sample-archive
if [ $? != 0 ]; then exit 1; fi
go install github.com/fossul/fossul/src/plugins/storage/basic/container-basic
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/storage/container-basic.so github.com/fossul/fossul/src/plugins/storage/native/container-basic
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/app/mariadb.so github.com/fossul/fossul/src/plugins/app/native/mariadb
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/app/mariadb-dump.so github.com/fossul/fossul/src/plugins/app/native/mariadb-dump
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/app/postgres.so github.com/fossul/fossul/src/plugins/app/native/postgres
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/app/postgres-dump.so github.com/fossul/fossul/src/plugins/app/native/postgres-dump
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/app/mongo.so github.com/fossul/fossul/src/plugins/app/native/mongo
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/app/mongo-dump.so github.com/fossul/fossul/src/plugins/app/native/mongo-dump
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/archive/aws.so github.com/fossul/fossul/src/plugins/archive/native/aws
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/storage/ocs-gluster.so github.com/fossul/fossul/src/plugins/storage/native/ocs-gluster
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/storage/csi-ceph.so github.com/fossul/fossul/src/plugins/storage/native/csi-ceph
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $PLUGIN_DIR/storage/csi.so github.com/fossul/fossul/src/plugins/storage/native/csi
if [ $? != 0 ]; then exit 1; fi

echo "Building Services"
go install github.com/fossul/fossul/src/engine/server
if [ $? != 0 ]; then exit 1; fi
go install github.com/fossul/fossul/src/engine/app
if [ $? != 0 ]; then exit 1; fi
go install github.com/fossul/fossul/src/engine/storage
if [ $? != 0 ]; then exit 1; fi

echo "Moving plugins to $PLUGIN_DIR"
if [ ! -d "$PLUGIN_DIR/app" ]; then mkdir $PLUGIN_DIR/app; fi
mv $GOBIN/sample-app $PLUGIN_DIR/app
if [ $? != 0 ]; then exit 1; fi
if [ ! -d "$PLUGIN_DIR/storage" ]; then mkdir $PLUGIN_DIR/storage; fi
mv $GOBIN/sample-storage $PLUGIN_DIR/storage
if [ $? != 0 ]; then exit 1; fi
if [ ! -d "$PLUGIN_DIR/archive" ]; then mkdir $PLUGIN_DIR/archive; fi
mv $GOBIN/sample-archive $PLUGIN_DIR/archive
if [ $? != 0 ]; then exit 1; fi
mv $GOBIN/container-basic $PLUGIN_DIR/storage
if [ $? != 0 ]; then exit 1; fi

echo "Copying default configs"
if [ ! -z $GOBIN/metadata/configs/default ]; then
  mkdir -p $GOBIN/metadata/configs/default
  if [ $? != 0 ]; then exit 1; fi
fi

cp -r $GOPATH/src/github.com/fossul/fossul/src/cli/configs/default $GOBIN/metadata/configs/default
if [ $? != 0 ]; then exit 1; fi

echo "Copying startup scripts"
cp $GOPATH/src/github.com/fossul/fossul/scripts/fossul-server-startup.sh $GOBIN
if [ $? != 0 ]; then exit 1; fi
cp $GOPATH/src/github.com/fossul/fossul/scripts/fossul-app-startup.sh $GOBIN
if [ $? != 0 ]; then exit 1; fi
cp $GOPATH/src/github.com/fossul/fossul/scripts/fossul-storage-startup.sh $GOBIN
if [ $? != 0 ]; then exit 1; fi

echo "Build completed successfully"

