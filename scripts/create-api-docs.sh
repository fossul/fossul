#!/bin/sh
SWAG=/home/ktenzer/swag

echo "Creating Server API docs"
cd $GOPATH/src/github.com/fossul/fossul/src/engine/server;$SWAG init -g server.go
if [ $? != 0 ]; then exit 1; fi

echo "Creating App API docs"
cd $GOPATH/src/github.com/fossul/fossul/src/engine/app;$SWAG init -g app.go
if [ $? != 0 ]; then exit 1; fi

echo "Creating Storage API docs"
cd $GOPATH/src/github.com/fossul/fossul/src/engine/storage;$SWAG init -g storage.go
if [ $? != 0 ]; then exit 1; fi

echo "API docs generated successfully"


