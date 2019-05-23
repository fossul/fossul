#!/bin/bash
PORT="8000"
CONFIG_DIR="metadata/configs"
DATA_DIR="metadata/data"
USERNAME="admin"
PASSWORD="redhat123"
SERVER_HOSTNAME="fossil-server"
SERVER_PORT="8000"
APP_HOSTNAME="fossil-app"
APP_PORT="8001"
STORAGE_HOSTNAME="fossil-storage"
STORAGE_PORT="8002"

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

if [[ -z "${FOSSIL_SERVER_CLIENT_HOSTNAME}" ]]; then
    export FOSSIL_SERVER_CLIENT_HOSTNAME=$SERVER_HOSTNAME
fi  

if [[ -z "${FOSSIL_SERVER_CLIENT_PORT}" ]]; then
    export FOSSIL_SERVER_CLIENT_PORT=$SERVER_PORT
fi  

if [[ -z "${FOSSIL_APP_CLIENT_HOSTNAME}" ]]; then
    export FOSSIL_APP_CLIENT_HOSTNAME=$APP_HOSTNAME
fi  

if [[ -z "${FOSSIL_APP_CLIENT_PORT}" ]]; then
    export FOSSIL_APP_CLIENT_PORT=$APP_PORT
fi  

if [[ -z "${FOSSIL_STORAGE_CLIENT_HOSTNAME}" ]]; then
    export FOSSIL_STORAGE_CLIENT_HOSTNAME=$STORAGE_HOSTNAME
fi  

if [[ -z "${FOSSIL_STORAGE_CLIENT_PORT}" ]]; then
    export FOSSIL_STORAGE_CLIENT_PORT=$STORAGE_PORT
fi  

$GOBIN/server