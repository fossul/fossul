#!/bin/bash
PORT="8002"
PLUGIN_DIR="/app/plugins"
USERNAME="admin"
PASSWORD="redhat123"

if [[ -z "${FOSSIL_STORAGE_PORT}" ]]; then
    export FOSSIL_STORAGE_PORT=$PORT
fi    

if [[ -z "${FOSSIL_STORAGE_PLUGIN_DIR}" ]]; then
    export FOSSIL_STORAGE_PLUGIN_DIR=$PLUGIN_DIR
fi    

if [[ -z "${FOSSIL_USERNAME}" ]]; then
    export FOSSIL_USERNAME=$USERNAME
fi   

if [[ -z "${FOSSIL_PASSWORD}" ]]; then
    export FOSSIL_PASSWORD=$PASSWORD
fi   

$GOBIN/storage

