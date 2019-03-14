#!/bin/sh

PLUGIN_DIR="/home/ktenzer/plugins"

echo "Installing Dependencies"
go get github.com/pborman/getopt/v2
go get github.com/gorilla/mux
go get k8s.io/apimachinery/pkg/apis/meta/v1
go get k8s.io/client-go/kubernetes/typed/core/v1 
go get k8s.io/client-go/rest
go get github.com/BurntSushi/toml

echo "Building Shared Libraries"
go build engine/util
go build engine/client
go build engine/client/k8s
go build engine/plugins/pluginUtil

echo "Building Plugins"
go install engine/plugins/sample-app
go install engine/plugins/sample-storage
go install engine/plugins/basic-container-storage

echo "Building Services"
go install engine/workflow
go install engine/app
go install engine/storage

echo "Moving plugins to $PLUGIN_DIR"
mv $GOBIN/sample-app $PLUGIN_DIR
mv $GOBIN/sample-storage $PLUGIN_DIR
mv $GOBIN/basic-container-storage $PLUGIN_DIR

echo "Build completed successfully"

