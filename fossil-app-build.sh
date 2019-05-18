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
go install fossil/src/engine/plugins/app/basic/sample-app
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSIL_BUILD_PLUGIN_DIR/app/sample-app.so fossil/src/engine/plugins/app/native/sample-app
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSIL_BUILD_PLUGIN_DIR/app/mariadb.so fossil/src/engine/plugins/app/native/mariadb
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSIL_BUILD_PLUGIN_DIR/app/mariadb-dump.so fossil/src/engine/plugins/app/native/mariadb-dump
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSIL_BUILD_PLUGIN_DIR/app/postgres.so fossil/src/engine/plugins/app/native/postgres
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSIL_BUILD_PLUGIN_DIR/app/postgres-dump.so fossil/src/engine/plugins/app/native/postgres-dump
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSIL_BUILD_PLUGIN_DIR/app/mongo.so fossil/src/engine/plugins/app/native/mongo
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSIL_BUILD_PLUGIN_DIR/app/mongo-dump.so fossil/src/engine/plugins/app/native/mongo-dump
if [ $? != 0 ]; then exit 1; fi

echo "Building App Service"
go install fossil/src/engine/app
if [ $? != 0 ]; then exit 1; fi

echo "Moving plugins to $FOSSIL_BUILD_PLUGIN_DIR"
if [ ! -d "$FOSSIL_BUILD_PLUGIN_DIR/app" ]; then mkdir $FOSSIL_BUILD_PLUGIN_DIR/app; fi
mv $GOBIN/sample-app $FOSSIL_BUILD_PLUGIN_DIR/app
if [ $? != 0 ]; then exit 1; fi

echo "Copying startup script"
cp $GOPATH/src/fossil/fossil-app-startup.sh $GOBIN
if [ $? != 0 ]; then exit 1; fi

echo "App build completed successfully"

