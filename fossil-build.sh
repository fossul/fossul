#!/bin/sh

PLUGIN_DIR="/home/ktenzer/plugins"

echo "Installing Dependencies"
go get github.com/pborman/getopt/v2
go get github.com/gorilla/mux
go get k8s.io/apimachinery/pkg/apis/meta/v1
go get k8s.io/client-go/kubernetes/typed/core/v1 
go get k8s.io/client-go/rest
go get github.com/BurntSushi/toml

echo "Running Unit Tests"
go test engine/util
if [ $? != 0 ]; then exit 1; fi
go test engine/plugins/pluginUtil
if [ $? != 0 ]; then exit 1; fi

echo "Building Shared Libraries"
go build engine/util
if [ $? != 0 ]; then exit 1; fi
go build engine/client
if [ $? != 0 ]; then exit 1; fi
go build engine/client/k8s
if [ $? != 0 ]; then exit 1; fi
go build engine/plugins/pluginUtil
if [ $? != 0 ]; then exit 1; fi

echo "Building Plugins"
go install engine/plugins/app/sample-app
if [ $? != 0 ]; then exit 1; fi
go install engine/plugins/storage/sample-storage
if [ $? != 0 ]; then exit 1; fi
go install engine/plugins/archive/sample-archive
if [ $? != 0 ]; then exit 1; fi
go install engine/plugins/storage/container-basic
if [ $? != 0 ]; then exit 1; fi

echo "Building Services"
go install engine/server
if [ $? != 0 ]; then exit 1; fi
go install engine/app
if [ $? != 0 ]; then exit 1; fi
go install engine/storage
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
cp -r $GOBIN/fossil/src/cli/configs/default $GOBIN/configs
if [ $? != 0 ]; then exit 1; fi

echo "Build completed successfully"

