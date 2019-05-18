#!/bin/bash
PORT="8000"
CONFIG_DIR="configs"
DATA_DIR="data"
USERNAME="admin"
PASSWORD="redhat123"

if [[ -z "${FOSSIL_SERVER_SERVICE_PORT}" ]]; then
    export FOSSIL_SERVER_SERVICE_PORT=$PORT
fi    

if [[ -z "${FOSSIL_SERVER_CONFIG_DIR}" ]]; then
    export FOSSIL_SERVER_CONFIG_DIR=$CONFIG_DIR
fi   

if [[ -z "${FOSSIL_SERVER_DATA_DIR}" ]]; then
    export FOSSIL_SERVER_DATA_DIR=$DATA_DIR
fi    

if [[ -z "${FOSSIL_USERNAME}" ]]; then
    export FOSSIL_USERNAME=$USERNAME
fi   

if [[ -z "${FOSSIL_PASSWORD}" ]]; then
    export FOSSIL_PASSWORD=$PASSWORD
fi   

$GOBIN/server

