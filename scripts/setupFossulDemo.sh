#!/bin/sh
CLI="/home/ktenzer/fossul"
CONFIG_DIR="../src/cli/configs/example"

$CLI --profile mariadb --action addProfile
$CLI --profile mongo --action addProfile
$CLI --profile postgres --action addProfile
$CLI --profile postgres --config postgres --action addConfig --config-file $CONFIG_DIR/postgres/postgres/postgres.conf
$CLI --profile postgres --config postgres --action addPluginConfig --config-file $CONFIG_DIR/postgres/postgres/container-basic.so.conf --plugin container-basic.so
$CLI --profile postgres --config postgres --action addPluginConfig --config-file $CONFIG_DIR/postgres/postgres/postgres-dump.so.conf --plugin postgres-dump.so
$CLI --profile mariadb --config mariadb --action addConfig --config-file $CONFIG_DIR/mariadb/mariadb/mariadb.conf
$CLI --profile mariadb --config mariadb --action addPluginConfig --config-file $CONFIG_DIR/mariadb/mariadb/container-basic.so.conf --plugin container-basic.so
$CLI --profile mariadb --config mariadb --action addPluginConfig --config-file $CONFIG_DIR/mariadb/mariadb/mariadb-dump.so.conf --plugin mariadb-dump.so
$CLI --profile mongo --config mongo --action addConfig --config-file $CONFIG_DIR/mongo/mongo/mongo.conf
$CLI --profile mongo --config mongo --action addPluginConfig --config-file $CONFIG_DIR/mongo/mongo/container-basic.so.conf --plugin container-basic.so
$CLI --profile mongo --config mongo --action addPluginConfig --config-file $CONFIG_DIR/mongo/mongo/mongo-dump.so.conf --plugin mongo-dump.so
$CLI --profile test --config test --action addConfig --config-file $CONFIG_DIR/test/test/test.conf
$CLI --profile test --config test --action addPluginConfig --config-file $CONFIG_DIR/test/test/sample-app.conf --plugin sample-app
$CLI --profile test --config test --action addPluginConfig --config-file $CONFIG_DIR/test/test/sample-storage.conf --plugin sample-storage
$CLI --profile test --config test --action addPluginConfig --config-file $CONFIG_DIR/test/test/sample-archive.conf --plugin sample-archive 
$CLI --profile mongo --config mongo --action addSchedule --cron-schedule "0,15,30,45 * * * *" --policy daily
$CLI --profile mongo --config mongo --action addSchedule --cron-schedule "20 * * * *" --policy weekly
$CLI --profile mariadb --config mariadb --action addSchedule --cron-schedule "0,15,30,45 * * * *" --policy daily
$CLI --profile mariadb --config mariadb --action addSchedule --cron-schedule "20 * * * *" --policy weekly
$CLI --profile postgres --config postgres --action addSchedule --cron-schedule "0,15,30,45 * * * *" --policy daily
$CLI --profile postgres --config postgres --action addSchedule --cron-schedule "20 * * * *" --policy weekly
$CLI --profile test --config test --action addSchedule --cron-schedule "0,15,30,45 * * * *" --policy daily
$CLI --profile test --config test --action addSchedule --cron-schedule "20 * * * *" --policy weekly
