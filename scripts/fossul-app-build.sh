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
go install fossul/src/engine/plugins/app/basic/sample-app
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/app/sample-app.so fossul/src/engine/plugins/app/native/sample-app
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/app/mariadb.so fossul/src/engine/plugins/app/native/mariadb
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/app/mariadb-dump.so fossul/src/engine/plugins/app/native/mariadb-dump
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/app/postgres.so fossul/src/engine/plugins/app/native/postgres
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/app/postgres-dump.so fossul/src/engine/plugins/app/native/postgres-dump
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/app/mongo.so fossul/src/engine/plugins/app/native/mongo
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/app/mongo-dump.so fossul/src/engine/plugins/app/native/mongo-dump
if [ $? != 0 ]; then exit 1; fi

echo "Building App Service"
go install fossul/src/engine/app
if [ $? != 0 ]; then exit 1; fi

echo "Moving plugins to $FOSSUL_BUILD_PLUGIN_DIR"
if [ ! -d "$FOSSUL_BUILD_PLUGIN_DIR/app" ]; then mkdir $FOSSUL_BUILD_PLUGIN_DIR/app; fi
mv $GOBIN/sample-app $FOSSUL_BUILD_PLUGIN_DIR/app
if [ $? != 0 ]; then exit 1; fi

echo "Copying startup script"
cp $GOPATH/src/fossul/scripts/fossul-app-startup.sh $GOBIN
if [ $? != 0 ]; then exit 1; fi

echo "App build completed successfully"

