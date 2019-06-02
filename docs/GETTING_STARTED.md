# Getting Started
Provides a how-to tutorial focused on Fossul operators.

## Getting Started using CLI
First fossul is deployed on OpenShift using the provided template in the yaml folder or K8s using Dockerfiles. You will end up with three pods, one for each service: server, storage and app.

### Download CLI
```curl https://raw.githubusercontent.com/ktenzer/fossul/master/release/fossul-cli_1.0.0.tar.gz |tar xz```

### Save Credentials
By default credentials will be stored under user home directory in a file called .fossul-credentials. You can specify `--credential-file` argument and path to save or load credential files from another location.

```$ fossul --set-credentials --user admin --pass redhat123 --server-host fossul-server-fossul.apps.46.4.207.247.xip.io --server-port 80 --app-host fossul-app-fossul.apps.46.4.207.247.xip.io --app-port 80 --storage-host fossul-storage-fossul.apps.46.4.207.247.xip.io --storage-port 80```

### Create a Profile
A profile is simply a group it can contain one or more configurations. For example you may have an application with several databases. The application is the profile and each database a configuration within the profile. It is only there for organizational purposes.
```$ fossul --profile mariadb --action addProfile```
 
### Create a COnfiguration
A configuration requires a main configuration and a configuration for each plugin used. A configuration just contains key/value pairs. The first step is to get the default configurations, change them locally and the upload them. Here we will create configuration to backup and restore mariadb. We will use container-basic and mariadb-dump plugins.
 
#### Main config
Simply copy past to file and update. ProfileName, ConfigName, WorkflowId, SelectedBackupPolicy, SelectedBackupRetention, SelectedWorkflowId are all ignored. These are added dynamically. All you need to do is add plugins app,storage, archive, configure auto discovery (depending on if plugin supports it), configure policy and any pre/post commands that should execute.

```$ fossul --get-default-config``` 

```
ProfileName = ""
ConfigName = ""
WorkflowId = ""
AppPlugin = "mariadb-dump.so"
StoragePlugin = "container-basic.so"
ArchivePlugin = ""
AutoDiscovery = true
JobRetention = 100
SelectedBackupPolicy = ""
SelectedBackupRetention = 0
SelectedWorkflowId = 0
PreAppQuiesceCmd = "echo,pre app quiesce command"
AppQuiesceCmd = "echo,app quiesce command"
PostAppQuiesceCmd = "echo,post app quiesce command"
BackupCreateCmd = "echo,backup create cmd"
BackupDeleteCmd = "echo,backup delete cmd"
ArchiveCreateCmd = "echo,archive create cmd"
ArchiveDeleteCmd = "echo,archive delete cmd"
PreAppUnquiesceCmd = "echo,pre app unquiesce command"
AppUnquiesceCmd = "echo,app unquiesce command"
PostAppUnquiesceCmd = "echo,post app unquiesce command"
PreAppRestoreCmd = ""
RestoreCmd = ""
PostAppRestoreCmd = ""
SendTrapErrorCmd = "echo,send trap error command"
SendTrapSuccessCmd = "echo,send trap success command"

[[BackupRetentions]]
  Policy = "daily"
  RetentionDays = 5

[[BackupRetentions]]
  Policy = "weekly"
  RetentionDays = 4
```
  
Assuming we saved file to /tmp/mariadb.conf
```$ fossul --profile mariadb --config mariadb --action addConfig --config-file /tmp/mariadb.conf```

#### Get a Config
```$ fossul --profile mariadb --config mariadb --get-config```

#### App Plugin Config
Almost same as previous step, get the plugin default config, copy/paste to file, file it out and add it back to server as new config.

```$ fossul --get-default-plugin-config --plugin mariadb-dump.so```

```
AccessWithinCluster = "true"
ContainerName = "mariadb"
MysqlDb = "sampledb"
MysqlDumpCmd = "/opt/rh/rh-mariadb102/root/usr/bin/mysqldump"
MysqlDumpPath = "/tmp"
MysqlHost = "localhost"
MysqlPort = "3306"
MysqlProto = "tcp"
MysqlRestoreCmd = "/opt/rh/rh-mariadb102/root/usr/bin/mysql"
MysqlUser = "root"
Namespace = "databases"
ServiceName = "mariadb"
```

Assuming we saved file to /tmp/mariadb-dump.conf
```$ fossul --profile mariadb --config mariadb --action addPluginConfig --plugin mariadb-dump.so --config-file /tmp/mariadb-dump.conf```

#### Get Plugin Config

```$ fossul --profile mariadb --config mariadb --get-plugin-config --plugin mariadb-dump.so```

#### Storage Plugin Configuration
Identical as previous step just a different plugin. The BackupSrcPaths option can be ignored if you set and your plugin supports auto-discover. The app plugin will automatically figure out what it should backup and set this dynamically.

```$ fossul --get-default-plugin-config --plugin container-basic.so```

```
AccessWithinCluster = "true"
BackupDestPath = "/app/backups"
BackupName = "cmds"
BackupSrcPaths = "/var/lib/mysql/data/sampledb,/var/lib/mysql/data/test"
ContainerPlatform = "openshift"
CopyCmdPath = "/app/oc"
Namespace = "databases"
ServiceName = "mariadb"
```

### Manual Backup
```$ fossul --profile mariadb --config mariadb --policy daily --action backup```

### List Backups
```$ fossul --profile mariadb --config mariadb --policy daily --action backupList```

### Add Schedule
Create a schedule that will run a backup every minute.
```$ fossul --profile mariadb --config mariadb --policy daily --action addSchedule --cron-schedule "* * * * *"```

### List Schedules
```$ fossul --list-schedules
### Job Schedules ###`
CronSchedule      ProfileName      ConfigName      Policy     
* * * * *         mariadb          mariadb         daily      
25 * * * *        mariadb          mariadb         weekly     
* * * * *         mongo            mongo           daily      
20 * * * *        mongo            mongo           weekly     
* * * * *         postgres         postgres        daily     
30 * * * *        postgres         postgres        weekly
```

### List jobs
```fossul --profile mariadb --config mariadb --action jobList 
########## Welcome to Fossul Framework ##########
### List of Jobs for profile [mariadb] config [mariadb] ###
WorkflowId      Type        Status        Policy      Start Time               
5886            backup      COMPLETE      daily       2019-05-23T23:29:00Z     
6777            backup      COMPLETE      daily       2019-05-23T23:28:00Z     
2652            backup      COMPLETE      daily       2019-05-23T23:27:00Z     
5338            backup      COMPLETE      daily       2019-05-23T23:26:00Z     
4967            backup      COMPLETE      weekly      2019-05-23T23:25:00Z
```
### Job Status
```$ fossul --profile mariadb --config mariadb --action jobStatus --workflow-id 6777```

### Restore
```$ fossul --profile mariadb --config mariadb --action restore --workflow-id 6777```