![](images/fossul_logo.png)
# Plugins
Purpose is to provide more details on plugins. Fossul is all about plugins. There are three types of plugins: application, storage and archive.

## Application Plugins
Application plugins are responsible for handling quiesce/unquiesce for backup or pre/post recovery operations.

### Sample-App
A sample plugin basic and native to provide example for building other plugins.

### MariaDB / MySQL
Both mariadb and mysql are supported in the same plugins.

#### Mariadb-Dump
This plugin will backup and recover mariadb or mysql databases using a dump, a logical backup.
The only difference between mysql or mariadb configuration is the path to the dump and restore commands.
For MariaDB:
* MysqlDumpCmd = "/opt/rh/rh-mariadb102/root/usr/bin/mysqldump"
* MysqlRestoreCmd = "/opt/rh/rh-mariadb102/root/usr/bin/mysql"
For MySQL
* MysqlDumpCmd = "/opt/rh/rh-mysql57/root/usr/bin/mysqldump"
* MysqlRestoreCmd = "/opt/rh/rh-mysql57/root/usr/bin/mysql"

#### Mariadb
This plugin does not use a dump and will quiesce/unquiesce the database in order to do a physical backup. It is recommended this plugin only be used when combined with snapshot technology. The quiesce will pause writes backup needs to happen in seconds.

### PostgreSQL

#### PostgreSQL-Dump
This plugin will backup and recover postgresql using dump, a logical backup.

#### PostgreSQL
This plugin does not use a dump and will quiesce/unquiesce the database in order to do a physical backup. It is recommended this plugin be used when combined with snapshot technology. Writes are not paused like with MySQL but you don't want to leave database in backup mode for extended time either. 

This plugin requires WAL archive logging to be configured in order to perform backups. This is not enabled by default under OpenShift. First create an archive directory under `/var/lib/pgsql/data/userdata` by connecting to the pod via rsh. Next update the following parameters in the `/var/lib/pgsql/data/userdata/postgresql.conf`.

* wal_level=archive
* archive_mode=on
* max_wal_senders=3
* archive_command = '/bin/cp %p /var/lib/pgsql/data/userdata/archive/%f'

### Mongo

#### Mongo-Dump 
This plugin will backup and recover mongo using dump, a logical backup.

#### Mongo
This plugin does not use a dump and will quiesce/unquiesce the database in order to do a physical backup. It is recommended this plugin only be used when combined with snapshot technology. The quiesce will pause writes backup needs to happen in seconds.

## Storage Plugins
Storage plugins are responsible for storage operations such as the physical backup and restore of data. Storage plugins integrate with vendor technology such as snapshots and expose it to the framework.

### Sample-Storage
A sample plugin basic and native to provide example for building other plugins.

### Container-Basic
The container basic plugin is a standard storage plugin that works regardless of storage used. It uses rsync in order to backup data from the container running the application to the container running the storage service. It supports both OpenShift and Kubernetes platforms.

## Archive Plugins
Archive plugins are responsible for archive operations such as archiving backups and recovering from archived backups using technologies such as S3.

### Sample-Archive
A sample plugin basic and native to provide example for building other plugins.