#!/bin/sh

echo "Setting CLI credentials"

$GOBIN/cli --set-credentials --user admin --pass $FOSSUL_PASSWORD --server-host fossul-server --server-port 8000 --app-host fossul-app --app-port 8001 --storage-host fossul-storage --storage-port 8002
if [ $? != 0 ]; then exit 1; fi

echo "CLI credentials successfully"

