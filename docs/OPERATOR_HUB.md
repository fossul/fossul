# OperatorHub Installation
Operator hub is like an app store for operators and k8s native applications. It ships with OpenShift but needs to be installed on other k8s distros. Operators including the fossul operator are available at [OperatorHub.io](https://operatorhub.io)

First you need to get cluster credentials. The below example is for gcloud k8s.
<pre>
$ gcloud container clusters get-credentials fossul-cluster --region us-central1-c
</pre>

Create the OperatorHub CRDs and k8s APIs hooks.
<pre>
$ kubectl create -f https://raw.githubusercontent.com/operator-framework/operator-lifecycle-manager/master/deploy/upstream/quickstart/crds.yaml
</pre>

Add cluster role binding. This is not always required but if you get permissions on next step it is needed.
<pre>
$ kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user myuser@gmail.com
</pre>

Deploy the OperatorHub operator and lifecycle manager.
<pre>
$ kubectl create -f https://raw.githubusercontent.com/operator-framework/operator-lifecycle-manager/master/deploy/upstream/quickstart/olm.yaml
</pre>

Deploy the Fossul operator from OperatorHub
<pre>
$ kubectl create -f https://operatorhub.io/install/fossul-operator.yaml -n operators
</pre>

Validate the Fossul operator is running
<pre>
$ kubectl get pods -n operators
NAME                               READY   STATUS    RESTARTS   AGE
fossul-operator-55896f5659-zkc49   1/1     Running   0          19h
</pre>
