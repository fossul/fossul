#!/bin/sh

echo "Creating Server API docs"
cd $GOPATH/src/fossul/src/engine/server;swag init -g server.go
if [ $? != 0 ]; then exit 1; fi

echo "Creating App API docs"
cd $GOPATH/src/fossul/src/engine/app;swag init -g app.go
if [ $? != 0 ]; then exit 1; fi

echo "Creating Storage API docs"
cd $GOPATH/src/fossul/src/engine/storage;swag init -g storage.go
if [ $? != 0 ]; then exit 1; fi

echo "API docs generated successfully"


