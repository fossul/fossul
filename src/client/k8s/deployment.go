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
	apps "k8s.io/api/apps/v1"
	//autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsclient "k8s.io/client-go/kubernetes/typed/apps/v1"
)

func GetDeployment(namespace, deploymentName, accessWithinCluster string) (*apps.Deployment, error) {
	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := appsclient.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	deploymentConfig, err := clientset.Deployments(namespace).Get(deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return deploymentConfig, nil
}

func GetDeploymentScaleInteger(namespace, deploymentName, accessWithinCluster string) (*int32, error) {
	deployment, err := GetDeployment(namespace, deploymentName, accessWithinCluster)
	if err != nil {
		return nil, err
	}

	return deployment.Spec.Replicas, nil
}

func ScaleDeployment(namespace, deploymentName, accessWithinCluster string, size int32) error {
	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return err
	}

	// create the clientset
	clientset, err := appsclient.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	deployment, err := GetDeployment(namespace, deploymentName, accessWithinCluster)
	if err != nil {
		return err
	}

	deployment.Spec.Replicas = &size

	_, err = clientset.Deployments(namespace).Update(deployment)
	if err != nil {
		return err
	}

	return nil
}
