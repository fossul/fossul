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
go install github.com/fossul/fossul/src/plugins/app/basic/sample-app
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/app/sample-app.so github.com/fossul/fossul/src/plugins/app/native/sample-app
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/app/mariadb.so github.com/fossul/fossul/src/plugins/app/native/mariadb
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/app/mariadb-dump.so github.com/fossul/fossul/src/plugins/app/native/mariadb-dump
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/app/postgres.so github.com/fossul/fossul/src/plugins/app/native/postgres
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/app/postgres-dump.so github.com/fossul/fossul/src/plugins/app/native/postgres-dump
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/app/mongo.so github.com/fossul/fossul/src/plugins/app/native/mongo
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $FOSSUL_BUILD_PLUGIN_DIR/app/mongo-dump.so github.com/fossul/fossul/src/plugins/app/native/mongo-dump
if [ $? != 0 ]; then exit 1; fi

echo "Building App Service"
go install github.com/fossul/fossul/src/engine/app
if [ $? != 0 ]; then exit 1; fi

echo "Moving plugins to $FOSSUL_BUILD_PLUGIN_DIR"
if [ ! -d "$FOSSUL_BUILD_PLUGIN_DIR/app" ]; then mkdir $FOSSUL_BUILD_PLUGIN_DIR/app; fi
mv $GOBIN/sample-app $FOSSUL_BUILD_PLUGIN_DIR/app
if [ $? != 0 ]; then exit 1; fi

echo "App build completed successfully"

