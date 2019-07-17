![](../images/fossul_logo.png)
# Deploy
Provides information on the various ways Fossul Framework can be deployed.

## Deploying Fossul on OpenShift
To make the installation robust and easy, templates are provided that are parameterized to deploy Fossul on OpenShift. Look at templates to see full list of parameters.
The deployment will first perform a build and then deploy. Grabbing a release branch instead of master which may be unstable is a good idea.

### Clone Github repository
```$ git clone https://github.com/fossul/fossul.git```

### Create Project
```$ oc new-project fossul```

### Add Security Context and CLuster Reader Permissions to Project
Fossul needs cluster-reader permissions and a security context that allows any user ID instead of the default range.

```$ oc process -f yaml/fossul-framework-cluster-permissions.yaml -p APPLICATION_NAME=fossul -p APPLICATION_NAMESPACE=fossul |oc create -f -```

Optionally you can also give the fossil project default service account cluster-admin permissions. This is not recommended but it saves you from adding access above and projects later.

```$ oc adm policy add-cluster-role-to-user cluster-admin system:serviceaccount:fossul:default```

### Join projects to allow network connectivity (optional) ###
Certain plugins access databases using the service.  you are using OpenShift SDN you need to allow network connectivity between projects where data workloads are running and fossul. In this case we have databases running in the 'databases' project.

```$ oc adm pod-network join-projects --to=databases fossul```

### Run Fossul Template
To deploy fossul by building source code use following template:

```$ oc process -f yaml/fossul-framework-template.yaml -p APPLICATION_NAME=fossul -p FOSSUL_USERNAME=admin -p FOSSUL_PASSWORD=r3dH@t31 |oc create -f -```

To deploy fossul using latest release without building the code:
```$ oc process -f yaml/fossul-framework-release-template.yaml -p APPLICATION_NAME=fossul -p FOSSUL_USERNAME=admin -p FOSSUL_PASSWORD=r3dH@t31 |oc create -f -```

### Add admin access to projects Fossul should manage
In order to run commands inside pods and containers, Fossul requires admin permissions to the project. In this example we have a project called databases and will add admin access from Fossul
default service account in order to manage pods running under the databases project.

```$ oc process -f yaml/fossul-framework-add-admin-access-to-project.yaml -p APPLICATION_NAME=fossul -p APPLICATION_NAMESPACE=fossul -p PROJECT_NAMESPACE=databases |oc create -f -```

### Remove Fossul OpenShift Deployment
To remove fossul completely first delete the project.

```$ oc delete project fossul```

Next remove role bindings and security context. First delete role binding in projects allowing admin access from Fossul project. In this case project called databases.

```$ oc delete rolebinding fossul-admin -n databases```

Remove the cluster rolebinding for reader access

```$ oc delete clusterrolebinding fossul-fossul-cluster-reader```

Remove security context allowing pods to run under project with any UID instead of range.

```$ oc delete scc fossul-scc```
