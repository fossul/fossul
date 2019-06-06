#!/bin/bash
PORT="8000"
CONFIG_DIR="metadata/configs"
DATA_DIR="metadata/data"
USERNAME="admin"
PASSWORD="redhat123"
SERVER_HOSTNAME="fossul-server"
SERVER_PORT="8000"
APP_HOSTNAME="fossul-app"
APP_PORT="8001"
STORAGE_HOSTNAME="fossul-storage"
STORAGE_PORT="8002"
DEBUG="false"

if [[ -z "${FOSSUL_SERVER_SERVICE_PORT}" ]]; then
    export FOSSUL_SERVER_SERVICE_PORT=$PORT
fi    

if [[ -z "${FOSSUL_SERVER_CONFIG_DIR}" ]]; then
    export FOSSUL_SERVER_CONFIG_DIR=$CONFIG_DIR
fi   

if [[ -z "${FOSSUL_SERVER_DATA_DIR}" ]]; then
    export FOSSUL_SERVER_DATA_DIR=$DATA_DIR
fi    

if [[ -z "${FOSSUL_USERNAME}" ]]; then
    export FOSSUL_USERNAME=$USERNAME
fi   

if [[ -z "${FOSSUL_PASSWORD}" ]]; then
    export FOSSUL_PASSWORD=$PASSWORD
fi  

if [[ -z "${FOSSUL_SERVER_CLIENT_HOSTNAME}" ]]; then
    export FOSSUL_SERVER_CLIENT_HOSTNAME=$SERVER_HOSTNAME
fi  

if [[ -z "${FOSSUL_SERVER_CLIENT_PORT}" ]]; then
    export FOSSUL_SERVER_CLIENT_PORT=$SERVER_PORT
fi  

if [[ -z "${FOSSUL_APP_CLIENT_HOSTNAME}" ]]; then
    export FOSSUL_APP_CLIENT_HOSTNAME=$APP_HOSTNAME
fi  

if [[ -z "${FOSSUL_APP_CLIENT_PORT}" ]]; then
    export FOSSUL_APP_CLIENT_PORT=$APP_PORT
fi  

if [[ -z "${FOSSUL_STORAGE_CLIENT_HOSTNAME}" ]]; then
    export FOSSUL_STORAGE_CLIENT_HOSTNAME=$STORAGE_HOSTNAME
fi  

if [[ -z "${FOSSUL_STORAGE_CLIENT_PORT}" ]]; then
    export FOSSUL_STORAGE_CLIENT_PORT=$STORAGE_PORT
fi  

if [[ -z "${FOSSUL_SERVER_DEBUG}" ]]; then
    export FOSSUL_SERVER_DEBUG=$DEBUG
fi

if [ ! -d "${FOSSUL_SERVER_CONFIG_DIR}/default" ]; then
    mkdir -p ${FOSSUL_SERVER_CONFIG_DIR}/default
    curl https://raw.githubusercontent.com/fossul/fossul/master/release/default-configs_1.0.0.tar.gz |tar xz;mv default ${FOSSUL_SERVER_CONFIG_DIR}/default
fi

$GOBIN/server