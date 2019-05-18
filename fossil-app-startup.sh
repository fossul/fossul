#!/bin/bash
PORT="8001"
PLUGIN_DIR="plugins"
USERNAME="admin"
PASSWORD="redhat123"

if [[ -z "${FOSSIL_APP_PORT}" ]]; then
    export FOSSIL_APP_PORT=$PORT
fi    

if [[ -z "${FOSSIL_APP_PLUGIN_DIR}" ]]; then
    export FOSSIL_APP_PLUGIN_DIR=$PLUGIN_DIR
fi    

if [[ -z "${FOSSIL_USERNAME}" ]]; then
    export FOSSIL_USERNAME=$USERNAME
fi   

if [[ -z "${FOSSIL_PASSWORD}" ]]; then
    export FOSSIL_PASSWORD=$PASSWORD
fi   

$GOBIN/app

