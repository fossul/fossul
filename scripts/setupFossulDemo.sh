#!/bin/sh
CLI="/home/ktenzer/fossul"
CREDENTIAL_FILE="/home/ktenzer/.fossul-credentials-remote"
CONFIG_DIR="/home/ktenzer/configs-remote/configs"

$CLI --profile mariadb --action addProfile --credential-file $CREDENTIAL_FILE 
$CLI --profile mongo --action addProfile --credential-file $CREDENTIAL_FILE 
$CLI --profile postgres --action addProfile --credential-file $CREDENTIAL_FILE 
$CLI --profile postgres --config postgres --action addConfig --config-file $CONFIG_DIR/postgres/postgres/postgres.conf --credential-file $CREDENTIAL_FILE 
$CLI --profile postgres --config postgres --action addPluginConfig --config-file $CONFIG_DIR/postgres/postgres/container-basic.so.conf --plugin container-basic.so --credential-file $CREDENTIAL_FILE 
$CLI --profile postgres --config postgres --action addPluginConfig --config-file $CONFIG_DIR/postgres/postgres/postgres-dump.so.conf --plugin postgres-dump.so --credential-file $CREDENTIAL_FILE 
$CLI --profile postgres --config postgres --action addPluginConfig --config-file $CONFIG_DIR/postgres/postgres/aws.so.conf --plugin aws.so --credential-file $CREDENTIAL_FILE 
$CLI --profile mariadb --config mariadb --action addConfig --config-file $CONFIG_DIR/mariadb/mariadb/mariadb.conf --credential-file $CREDENTIAL_FILE 
$CLI --profile mariadb --config mariadb --action addPluginConfig --config-file $CONFIG_DIR/mariadb/mariadb/container-basic.so.conf --plugin container-basic.so --credential-file $CREDENTIAL_FILE 
$CLI --profile mariadb --config mariadb --action addPluginConfig --config-file $CONFIG_DIR/mariadb/mariadb/mariadb-dump.so.conf --plugin mariadb-dump.so --credential-file $CREDENTIAL_FILE 
$CLI --profile mariadb --config mariadb --action addPluginConfig --config-file $CONFIG_DIR/mariadb/mariadb/aws.so.conf --plugin aws.so --credential-file $CREDENTIAL_FILE 
$CLI --profile mongo --config mongo --action addConfig --config-file $CONFIG_DIR/mongo/mongo/mongo.conf --credential-file $CREDENTIAL_FILE 
$CLI --profile mongo --config mongo --action addPluginConfig --config-file $CONFIG_DIR/mongo/mongo/container-basic.so.conf --plugin container-basic.so --credential-file $CREDENTIAL_FILE 
$CLI --profile mongo --config mongo --action addPluginConfig --config-file $CONFIG_DIR/mongo/mongo/mongo-dump.so.conf --plugin mongo-dump.so --credential-file $CREDENTIAL_FILE 
$CLI --profile mongo --config mongo --action addPluginConfig --config-file $CONFIG_DIR/mongo/mongo/aws.so.conf --plugin aws.so --credential-file $CREDENTIAL_FILE 
$CLI --profile mongo --config mongo --action addSchedule --cron-schedule "0,15,30,45 * * * *" --credential-file $CREDENTIAL_FILE --policy daily
$CLI --profile mongo --config mongo --action addSchedule --cron-schedule "20 * * * *" --credential-file $CREDENTIAL_FILE --policy weekly
$CLI --profile mariadb --config mariadb --action addSchedule --cron-schedule "0,15,30,45 * * * *" --credential-file $CREDENTIAL_FILE --policy daily
$CLI --profile mariadb --config mariadb --action addSchedule --cron-schedule "20 * * * *" --credential-file $CREDENTIAL_FILE --policy weekly
$CLI --profile postgres --config postgres --action addSchedule --cron-schedule "0,15,30,45 * * * *" --credential-file $CREDENTIAL_FILE --policy daily
$CLI --profile postgres --config postgres --action addSchedule --cron-schedule "20 * * * *" --credential-file $CREDENTIAL_FILE --policy weekly
