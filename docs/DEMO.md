![](../images/fossul_logo.png)
# Demo Environment
A quick way to setup and demo environment with fossul once you have deployed fossul to a project or namespace called fossul.

## Clone Fossul repository
```$ git clone https://github.com/fossul/fossul.git```

## Create databases project
```$ oc new-project databases```

## Deploy mariadb
```$ oc process -n openshift mariadb-persistent | oc create -f -```

## Deploy postgres
```$ oc process -n openshift postgresql-persistent | oc create -f -```

## Deploy mongo
```$ oc process -n openshift mongodb-persistent | oc create -f -```

## Allow fossul access to database project
From database project run the template to add access to fossul project.

```$ oc process -f yaml/fossul-framework-add-admin-access-to-project.yaml -p PROJECT_NAMESPACE=fossul |oc create -f -```

## Extract cli tarball
```$ tar -xvf release/fossul-cli_linux_1.0.0.tar.gz```

## Save fossul credentials file
You need the route to the server, storage and app services. Port should be 80 and user/pass is whatever you used when deploying template.

```$ fossul --set-credentials --user admin --pass redhat123 --server-host fossul-server-fossul.apps.46.4.207.247.xip.io --server-port 80 --app-host fossul-app-fossul.apps.46.4.207.247.xip.io --app-port 80 --storage-host fossul-storage-fossul.apps.46.4.207.247.xip.io --storage-port 80```

## Setup demo environment
A script is provided to setup configs for a demo environment. You need to update path to CLI and ensure path to example configs directory is correct. This will not only create configurations but a job schedule which will execute backups. You can see what commands were run to understand how to set things up in future on your own.

```$ cd scripts```

```$ ./setupFossulDemo.sh```

## Check job schedules
```$ fossul --list-schedules```
