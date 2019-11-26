/*
Copyright 2019 The Fossul Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package k8s

import (
	"errors"
	deploymentConfigClient "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	"k8s.io/client-go/kubernetes"
	deploymentClient "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

func getClient(accessWithinCluster string) (*kubernetes.Clientset, error) {
	var client *kubernetes.Clientset
	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return client, err
	}

	client, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return client, err
	}

	return client, nil
}

func getDeploymentConfigClient(accessWithinCluster string) (*deploymentConfigClient.AppsV1Client, error) {
	var client *deploymentConfigClient.AppsV1Client
	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return client, err
	}

	client, err = deploymentConfigClient.NewForConfig(kubeConfig)
	if err != nil {
		return client, err
	}

	return client, nil
}

func getDeploymentClient(accessWithinCluster string) (*deploymentClient.AppsV1Client, error) {
	var client *deploymentClient.AppsV1Client
	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return client, err
	}

	client, err = deploymentClient.NewForConfig(kubeConfig)
	if err != nil {
		return client, err
	}

	return client, nil
}

func getKubeConfig(accessWithinCluster string) (error, *rest.Config) {
	var kubeConfig *rest.Config
	var err error
	if accessWithinCluster == "true" {
		kubeConfig, err = rest.InClusterConfig()
		if err != nil {
			return err, nil
		}
	} else if accessWithinCluster == "false" {
		var kubeconfigFile string
		if home := homeDir(); home != "" {
			kubeconfigFile = home + "/.kube" + "/config"
			if _, err := os.Stat(kubeconfigFile); os.IsNotExist(err) {
				log.Println(err, "\n"+"[ERROR] Kube config not found under "+kubeconfigFile)
				return err, nil
			}
		} else {
			log.Println("[ERROR] Could not find homedir, check environment!")
			return err, nil
		}

		// use the current context in kubeconfig
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfigFile)
		if err != nil {
			log.Println(err.Error())
			return err, nil
		}
	} else {
		log.Println("[ERROR]: Parameter AccessWithinCluster not set to true or false")
		err := errors.New("Parameter AccessWithinCluster not set to true or false")
		return err, nil
	}

	return nil, kubeConfig
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func GetOwner(namespace, service, accessWithinCluster string) (string, string, error) {
	podName, err := GetPodName(namespace, service, accessWithinCluster)
	if err != nil {
		return "", "", err
	}

	pod, err := GetPod(podName, namespace, accessWithinCluster)
	if err != nil {
		return "", "", err
	}

	var rcName string
	for _, owner := range pod.OwnerReferences {
		rcName = owner.Name
	}

	rc, err := GetReplicationController(rcName, namespace, accessWithinCluster)
	if err != nil {
		return "", "", err
	}

	var ownerKind string
	var ownerName string
	for _, owner := range rc.OwnerReferences {
		ownerKind = owner.Kind
		ownerName = owner.Name
	}

	return ownerKind, ownerName, nil
}
