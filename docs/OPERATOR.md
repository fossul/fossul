# Installation
It is recommended to install the operator into the openshift-operators namespace which is the default location for operators. Once the operator is installed create a new project called ```fossul``` and create the Fossul custom resource under the fossul namespace. This will deploy the Fossul Framework. You can choose another namespace for the Fossul Framework but then also need to provide the optional 'fossul_namespace' parameter with the spec for all the additional custom resources. This is required so the Fossul operator can communicate with Fossul.

Once Fossul is deployed you can manage backups either through custom resources provided by the operator or the Fossul CLI/API. If using the customer resources, first create a BackupConfig custom resource for every application you want to backup. Once that is done, you can optionally create a backup by simply creating a Backup custom resource. You can also create a BackupSchedule custom resource using cron syntax which will schedule your backup to occur on a regular interval or do so via a Kubernetes job. Finally if a backup has been created you can perform a restore through a custom resource as well, providing the workflow id from the completed backup. If you are creating backups through the custom resource the workflow id will be appended to the spec once the backup completes. If you are using the Fossul CLI/API or the BackupSchedule custom resource to create backups you will need to get the workflow id through the Fossul CLI/API directly as the Fossul scheduler won't create Fossul Backup custom resources.

For users of the cli, a cli pod is deployed with credentials and if desired you can create or manage backups from within that pod. To use the API you need the credentials stored in the fossul secret. The cli has already been pre-configured with these credentials. The API can also be optionally exposed via routes and API documentation is under ```/api/v1/index.html```.

# Quick Setup
When creating Fossul custom resources generally you will only need to update the deployment_name, deployment_type (DeploymentConfig, Deployment or VirtualMachine) and namespace. All custom resources must be created in the namespace where the database or application exists. 

## Deploy Fossul in fossul namespace
```$ oc new-project fossul```
<pre>
$ vi fossul.yaml
kind: Fossul
apiVersion: fossul.io/v1
metadata:
  name: fossul-sample
  namespace: fossul
spec: {}
</pre>
```$ oc create -f fossul.yaml```

## Deploy Mariadb Database
Make sure you create all Fossul custom resources in the same namespace as the database, application or virtual machine

## Create MariaDB Fossul BackupConfig
Creates a backup configuration for MariaDB or MySQL databases

```$ vi backupconfig.yaml```
<pre>
kind: BackupConfig
apiVersion: fossul.io/v1
metadata:
  name: mariadb-sample
  namespace: databases
spec:
  auto_discovery: false
  deployment_name: mariadb
  deployment_type: DeploymentConfig
  job_retention: 50
  overwrite_pcv_on_restore: true
  policies:
  - policy: hourly
    retentionNumber: 3
  - policy: daily
    retentionNumber: 10
  pvc_deletion_timeout: 300
  restore_to_new_pvc: false
  snapshot_timeout: 180
  storage_plugin: csi.so
  app_plugin: mariadb.so 
</pre>
```$ oc create -f backupconfig.yaml```

## Create PostgreSQL Fossul BackupConfig
Creates a backup config for PostgreSQL databases, you need to ensure the user defined in secret has backup permissions

```$ vi backupconfig.yaml```
<pre>
kind: BackupConfig
apiVersion: fossul.io/v1
metadata:
  name: postgresql-sample
  namespace: databases
spec:
  auto_discovery: false
  deployment_name: postgresql 
  deployment_type: DeploymentConfig
  job_retention: 50
  overwrite_pcv_on_restore: true
  policies:
  - policy: hourly
    retentionNumber: 3
  - policy: daily
    retentionNumber: 10
  pvc_deletion_timeout: 300
  restore_to_new_pvc: false
  snapshot_timeout: 180
  storage_plugin: csi.so
  app_plugin: postgres.so 
</pre>
```$ oc create -f backupconfig.yaml```

## Create Kubevirt Fossul BackupConfig
Creates a backup configuration for virtual machines

```$ vi backupconfig.yaml```
<pre>
kind: BackupConfig
apiVersion: fossul.io/v1
metadata:
  name: rhel8-vm-sample
  namespace: virtualmachines
spec:
  auto_discovery: false
  deployment_name: rhel8-vm
  deployment_type: VirtualMachine
  job_retention: 50
  overwrite_pcv_on_restore: true
  policies:
  - policy: hourly
    retentionNumber: 3
  - policy: daily
    retentionNumber: 10
  pvc_deletion_timeout: 300
  restore_to_new_pvc: false
  snapshot_timeout: 180
  storage_plugin: csi.so
  app_plugin: kubevirt.so 
</pre>
```$ oc create -f backupconfig.yaml```

## Create BackupSchedule
A backup can be scheduled per policy, defined by backup configuration using cron syntax

```$ vi backupschedule.yaml```
<pre>
kind: BackupSchedule
apiVersion: fossul.io/v1
metadata:
  name: mariadb-sample
  namespace: databases
spec:
  cron_schedule: '59 23 * * *'
  deployment_name: mariadb
  policy: daily
</pre>
```$ oc create -f backupschedule.yaml```

## Create Backup
A backup will be created according to policy and deleted either manually or as defined in retention policy

```$ vi backup.yaml```
<pre>
kind: Backup
apiVersion: fossul.io/v1
metadata:
  name: mariadb-sample
  namespace: databases
spec:
  deployment_name: mariadb
  policy: daily
</pre>
```$ oc create -f backup.yaml```

## Perform Restore
A restore requires a workflow_id located in the backup spec, upon a successful restore the custom resource will deleted automatically

```$ vi restore.yaml```
<pre>
kind: Restore
apiVersion: fossul.io/v1
metadata:   
  name: mariadb-sample
  namespace: databases
spec:
  deployment_name: mariadb
  policy: daily
  workflow_id: xxxx
  </pre>
```$ oc create -f restore.yaml```