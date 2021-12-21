#!/bin/bash

if [[ -z "${APP_PLUGIN_DIR}" ]]; then
    APP_DIR="${HOME}/plugins/app"
    export APP_PLUGIN_DIR=$APP_DIR

  if [[ ! -e "${APP_PLUGIN_DIR}" ]]; then
      mkdir -p  $APP_PLUGIN_DIR
  fi
fi

if [[ "$CI" == "true" ]]; then
    export APP_PLUGIN_DIR="."
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
go build -buildmode=plugin -o $APP_PLUGIN_DIR/sample-app.so github.com/fossul/fossul/src/plugins/app/native/sample-app
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $APP_PLUGIN_DIR/mariadb.so github.com/fossul/fossul/src/plugins/app/native/mariadb
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $APP_PLUGIN_DIR/mariadb-dump.so github.com/fossul/fossul/src/plugins/app/native/mariadb-dump
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $APP_PLUGIN_DIR/postgres.so github.com/fossul/fossul/src/plugins/app/native/postgres
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $APP_PLUGIN_DIR/postgres-dump.so github.com/fossul/fossul/src/plugins/app/native/postgres-dump
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $APP_PLUGIN_DIR/mongo.so github.com/fossul/fossul/src/plugins/app/native/mongo
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $APP_PLUGIN_DIR/mongo-dump.so github.com/fossul/fossul/src/plugins/app/native/mongo-dump
if [ $? != 0 ]; then exit 1; fi
go build -buildmode=plugin -o $APP_PLUGIN_DIR/kubevirt.so github.com/fossul/fossul/src/plugins/app/native/kubevirt
if [ $? != 0 ]; then exit 1; fi

echo "Building App Service"
go install github.com/fossul/fossul/src/engine/app
if [ $? != 0 ]; then exit 1; fi

if [[ ! -z "${GOBIN}" ]]; then
	echo "Moving plugins to $APP_PLUGIN_DIR"
	mv $GOBIN/sample-app $APP_PLUGIN_DIR

	echo "Copying startup scripts"
	cp scripts/fossul-app-startup.sh $GOBIN
	if [ $? != 0 ]; then exit 1; fi
fi


echo "App build completed successfully"
