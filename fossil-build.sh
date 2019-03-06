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

echo "Building Plugins"
go install engine/app/plugins/sample-app
go install engine/app/plugins/sample-storage

echo "Building Services"
go install engine/workflow
go install engine/app
go install engine/storage

echo "Moving plugins to $PLUGIN_DIR"
mv $GOBIN/sample-app $PLUGIN_DIR
mv $GOBIN/sample-storage $PLUGIN_DIR

echo "Build completed successfully"

