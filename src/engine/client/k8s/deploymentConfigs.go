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
	apps "github.com/openshift/api/apps/v1"
	appsclient "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetDeploymentConfig(namespace, deploymentConfigName, accessWithinCluster string) (*apps.DeploymentConfig, error) {
	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := appsclient.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	deploymentConfig, err := clientset.DeploymentConfigs(namespace).Get(deploymentConfigName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return deploymentConfig, nil
}

func GetDeploymentConfigScaleInteger(namespace, deploymentConfigName, accessWithinCluster string) (int32, error) {
	deploymentConfig, err := GetDeploymentConfig(namespace, deploymentConfigName, accessWithinCluster)
	if err != nil {
		return 0, err
	}

	return deploymentConfig.Spec.Replicas, nil
}

func ScaleDeploymentConfig(namespace, deploymentConfigName, accessWithinCluster string, size int32) error {
	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return err
	}

	// create the clientset
	clientset, err := appsclient.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	deploymentConfig, err := GetDeploymentConfig(namespace, deploymentConfigName, accessWithinCluster)
	if err != nil {
		return err
	}

	deploymentConfig.Spec.Replicas = size

	_, err = clientset.DeploymentConfigs(namespace).Update(deploymentConfig)
	if err != nil {
		return err
	}

	return nil
}
