# General
Fossil is a container-native backup and recovery framework. It aims to provide an application consistent backup, restore as well as recorvery for container-native applications and databases. It's goal is to provide enterprise backup and recovery capabilities enjoyed in traditional world, container-native, in order to increase adoption and make the migration from traditial to container that much easier. The challenge of backup, especially in the container-driven world is signifigant. Each application or database has it's own tools, methods and even procedures. Applications are much more fluid, dynamic and provide greater abstractions. In addition there is an enormous variation of application and storage vendors capabilities (for example snapshots that impact backup / restore). 

Fossil addresses the challanges by building a modular plugin based framework. Backup and recovery is always the same process, the only thing that changes is the specific tools and proceadures. The fossil belief is that the backup/recovery process can be standardized or democratized using a dynamic plugin-driven workflow and framework. Welcome to Fossil where you can not only imortalize your applications but their data as well!

In fossil there are three types of plugins storage, application and archive. 

## Storage
Storage plugins expose the capabilities of the underlying storage. They are responsible for the physical data and backing up as well as restoring it. Storage plugins in fossil run under the storage micro-service.

## Application
Application plugins expose the capabilities of the application. Before a backup is taken the application must be quiesced or dumped. Once data is restored application recovery must be performed to bring the application back into operation using a specific dataset. These operations are performed by an application plugin. Application plugins run under the application micro-service

## Archive
Archive plugins expose capabilities of secondary or tiertiary storage. Data is sacred so just having a single copy, on likely expensive storage is not always good enough. Archiving backups to something like S3 allows for longer-term storage at a cheaper cost. It also protects against losing the initial backup for whatever reason. Archive plugins are very close to storage and as such these plugins also run under the storage micro-service.

## Server
In fossil the server micro-service is an orchestration layer. This is where the workflow lives and is executed. There are two workflows one for backup and the other restore. Fossil could however easily support additional workflows in the future. The server handles communication with the other services, provides standard messaging, error handling and state (configurations, jobs, workflows, etc).

# Plugin Frameowrk
Fossil provides an extensive plugin framework. Plugins can be written in any language. There are two types of plugins native and basic.

## Native Plugins
Native plugins are written in go and loaded as a shared library. Native plugins also enforce plugin methods through an interface and are statically linked within the fossil plugin loader. Native plugins can however be written as well as compiled independently of the fossil framework. The fossil framework provides many client and plugin utility classes that are written in Go and native plugins get access to them. Native plugins also inherit their configuration as a object and can access and use API Objects without dealing with outputting JSON.

## Basic Plugins
Basic plugins are any plugins not written in Go or are not loaded using the fossil plugin loader. In other words that are not loaded as a shared library. Basic plugins can be written in any language, even shell script as they are executed via a system call. Their configuration is inherited via environment variables. For messages there is a standard output parser so a plugin simply needs to write to STDOUT in a specific format. Some methods, like pluginInfo (which exposes the methods, version and other information related to plugin) require JSON output.

There is really no advantage or disadvantage of native vs basic plugins. The framework provides plenty of examples of both. Native plugins require a static link in the plugin loader so the idea is that native plugins are added with new releases and basic plugins can be added anytime. Also as a plugin matures and becomes accepted in the framework it moves from a basic to native plugin. Basically it is up to you to decide.

# Contribution
If you want to provide additional features, please feel free to contribute via pull requests and please document your pull request thouroughly. We are happy to track and discuss ideas, topics and requests via 'Issues'. It is recommended you should start by writing plugins, as you can jump-in and write them with almost no learning curve.

# Architecture
The fossil architecture is a plugin based framework. Everything that happens in implemented in form of a method that executes inside a plugin. There are three main components that run as microservices inside containers. The server, application and storage services. As mentioned already the server is responsible for state, coordination and workflow. The application service is responsible for executing application plugins. The storage service is responsible for executing storage and archive plugins. Each service exposes it's own APIs. Nothing happens or can be done without issuing an API. In addition there is the CLI which is built on client libraries that marshall requests to the API. The framework, services, CLI and initial plugins are written in Go. Additional plugins can be written in any language as mentioned, you choose!

![](images/fossul_architecture_1.0.0.png)

# Workflow Engine
Workflows and the ability to democratoize a process like backup or restore is the key to fossil. In fossil a workflow has it's own Id and a series of steps. Each workflow step is an API to a plugin or CMD that executes the step. In fossil you could just use commands and not even any plugins. The plugins or commands which are executed are decided upon within a configuration. A fossil workflow takes as input a configuration. Configurations also define any pre/post commands (simple commands or scripts that can be executed in workflow) and also the backup policy as well as retention. Each plugin also has it's own configuration. These are all loaded and added to the config which is passed into all plugin operations or calls. As mentioned above in case of basic plugin the config object is demarshalled into environment variables. Every workflow has it's own log of what happened during workflow execution. You can decide in configuration how long to keep workflows. Finally a workflow has a state RUNNING, COMPLETE or ERROR.

# Profile
A profile is just an organizational unit or group of configurations.

# Configurations
A configuration lives under a profile. There are two types of configurations: main configuration and plugin configuration. A configuration is represented in TOML, a configuration file format. Configurations are pulled from the fossil server, edited in a file and then added back as a new configuration. This makes for maintaining and editing configurations very fast and easy.

# Job Scheduler
Fossil provides a job scheduler for scheduling of the various workflows. The scheduler implements a cron-style scheduler that utilizes cron syntax. Scheduler APIs are provided by the server service and scheduler job state is also stored on the server.

# Commands
Fossil framework allows user-defined commands to be executed via system calls. The main configuration has all the commands that can be executed. In fact you could just not use any plugins and do everything via commands if you wanted. The main idea though is to augment and provide maybe some special task capabilities that plugins aren't able to do through commands.

The commands itself is separated by its arguments via a ','. For example to print hello world as preQuiesceCmd you would do following.
```PreAppQuiesceCmd = "echo, hello world"```

# Plugins
Fossil is all about plugins. There are three types of plugins: application, storage and archive.

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

# CLI
The fossil CLI consumes APIs provided by the various services using client libraries that implement and marshall the various rest calls. The CLI can be used remotely and requires a credentials file, usually stored in user home directory to access the various services. 

# API
API documentation is done using swagger for Go (swaggo). Each services provides API documentation for the APIs that are published by the service.
* [http://fossil-server-fossil.apps.46.4.207.247.xip.io/api/v1/index.html](http://fossil-server-fossil.apps.46.4.207.247.xip.io/api/v1/index.html)
* [http://fossil-app-fossil.apps.46.4.207.247.xip.io/api/v1/index.html](http://fossil-app-fossil.apps.46.4.207.247.xip.io/api/v1/index.html)
* [http://fossil-storage-fossil.apps.46.4.207.247.xip.io/api/v1/index.html](http://fossil-storage-fossil.apps.46.4.207.247.xip.io/api/v1/index.html)

# Getting Started

## Development Environment
The instructions are for a Fedora 28 development environment but any Linux or MacOS should work.

### Install the Go programming language. 
```$ sudo dnf install -y go```

To build source code and setup a development ensure the following environment parameters are exported to the shell and set in user profile (.bashrc):

* export GOPATH=/home/fedora/go
* export GOBIN=/home/fedora

### Download dep binary. Dep is used for dependency and package management. Build scripts will call dep to download correct dependencies.
```$ curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh```

This will download and install dep into $GOBIN

### Clone the Fossil Github repository from '$GOPATH/src' in this case '/home/fedora/go/src'.
```$ git clone https://github.com/ktenzer/fossil.git```

### Change directory to the Fossil Github repository
```$ cd /home/fedora/go/src/fossil```

### Update Plugin Dir parameter in fossil build script
```
vi fossil-build.sh
PLUGIN_DIR="/home/fedora/plugins"
```

### Run the fossil build script
```$ /home/fedora/go/src/fossil/fossil-build.sh```

## Deploying Fossil

### OpenShift
An OpenShift template is provided under the yaml folder to deploy the fossil framework on OpenShift.

#### Clone Github repository
```$ git clone https://github.com/ktenzer/fossil.git```

#### Create Project
```$ oc create project fossil```

#### Add cluster permissions to fossil project
```$ oc adm policy add-cluster-role-to-user cluster-admin system:serviceaccount:fossil:default```

#### Deploy using template
```$ oc create -f fossil/yaml/fossil-engine-template.yaml -n fossil```

## Getting Started using CLI
First fossil is deployed on OpenShift using the provided template in the yaml folder or K8s using Dockerfiles. You will end up with three pods, one for each service: server, storage and app.

### Download CLI
```curl https://raw.githubusercontent.com/ktenzer/fossil/master/release/fossil-cli_1.0.0.tar.gz |tar xz```

### Save Credentials
By default credentials will be stored under user home directory in a file called .fossil-credentials. You can specify `--credential-file` argument and path to save or load credential files from another location.

```$ fossil --set-credentials --user admin --pass redhat123 --server-host fossil-server-fossil.apps.46.4.207.247.xip.io --server-port 80 --app-host fossil-app-fossil.apps.46.4.207.247.xip.io --app-port 80 --storage-host fossil-storage-fossil.apps.46.4.207.247.xip.io --storage-port 80```

### Create a Profile
A profile is simply a group it can contain one or more configurations. For example you may have an application with several databases. The application is the profile and each database a configuration within the profile. It is only there for organizational purposes.
```$ fossil --profile mariadb --action addProfile```
 
### Create a COnfiguration
A configuration requires a main configuration and a configuration for each plugin used. A configuration just contains key/value pairs. The first step is to get the default configurations, change them locally and the upload them. Here we will create configuration to backup and restore mariadb. We will use container-basic and mariadb-dump plugins.
 
#### Main config
Simply copy past to file and update. ProfileName, ConfigName, WorkflowId, SelectedBackupPolicy, SelectedBackupRetention, SelectedWorkflowId are all ignored. These are added dynamically. All you need to do is add plugins app,storage, archive, configure auto discovery (depending on if plugin supports it), configure policy and any pre/post commands that should execute.

```$ fossil --get-default-config``` 

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
```$ fossil --profile mariadb --config mariadb --action addConfig --config-file /tmp/mariadb.conf```

#### Get a Config
```$ fossil --profile mariadb --config mariadb --get-config```

#### App Plugin Config
Almost same as previous step, get the plugin default config, copy/paste to file, file it out and add it back to server as new config.

```$ fossil --get-default-plugin-config --plugin mariadb-dump.so```

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
```$ fossil --profile mariadb --config mariadb --action addPluginConfig --plugin mariadb-dump.so --config-file /tmp/mariadb-dump.conf```

#### Get Plugin Config

```$ fossil --profile mariadb --config mariadb --get-plugin-config --plugin mariadb-dump.so```

#### Storage Plugin Configuration
Identical as previous step just a different plugin. The BackupSrcPaths option can be ignored if you set and your plugin supports auto-discover. The app plugin will automatically figure out what it should backup and set this dynamically.

```$ fossil --get-default-plugin-config --plugin container-basic.so```

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
```$ fossil --profile mariadb --config mariadb --policy daily --action backup```

### List Backups
```$ fossil --profile mariadb --config mariadb --policy daily --action backupList```

### Add Schedule
Create a schedule that will run a backup every minute.
```$ fossil --profile mariadb --config mariadb --policy daily --action addSchedule --cron-schedule "* * * * *"```

### List Schedules
```$ fossil --list-schedules
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
```fossil --profile mariadb --config mariadb --action jobList 
########## Welcome to Fossil Framework ##########
### List of Jobs for profile [mariadb] config [mariadb] ###
WorkflowId      Type        Status        Policy      Start Time               
5886            backup      COMPLETE      daily       2019-05-23T23:29:00Z     
6777            backup      COMPLETE      daily       2019-05-23T23:28:00Z     
2652            backup      COMPLETE      daily       2019-05-23T23:27:00Z     
5338            backup      COMPLETE      daily       2019-05-23T23:26:00Z     
4967            backup      COMPLETE      weekly      2019-05-23T23:25:00Z
```
### Job Status
```$ fossil --profile mariadb --config mariadb --action jobStatus --workflow-id 6777```

### Restore
```$ fossil --profile mariadb --config mariadb --action restore --workflow-id 6777```
