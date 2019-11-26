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
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

func GetDeployment(namespace, deploymentName, accessWithinCluster string) (*apps.Deployment, error) {
	var deployment *apps.Deployment
	client, err := getDeploymentClient(accessWithinCluster)
	if err != nil {
		return deployment, err
	}

	deployment, err = client.Deployments(namespace).Get(deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func GetDeploymentScaleInteger(namespace, deploymentName, accessWithinCluster string) (*int32, error) {
	deployment, err := GetDeployment(namespace, deploymentName, accessWithinCluster)
	if err != nil {
		return nil, err
	}

	return deployment.Spec.Replicas, nil
}

func ScaleDownDeployment(namespace, deploymentConfigName, accessWithinCluster string, size int32, t int) error {
	client, err := getDeploymentClient(accessWithinCluster)
	if err != nil {
		return err
	}

	deployment, err := GetDeployment(namespace, deploymentConfigName, accessWithinCluster)
	if err != nil {
		return err
	}

	deployment.Spec.Replicas = &size

	_, err = client.Deployments(namespace).Update(deployment)
	if err != nil {
		return err
	}

	var poll = 5 * time.Second
	timeout := time.Duration(t) * time.Second
	start := time.Now()

	fmt.Printf("[DEBUG] Waiting up to %v for deployment to be scaled to %d\n", timeout, size)

	return wait.PollImmediate(poll, timeout, func() (bool, error) {
		deploymentConfig, err := GetDeployment(namespace, deploymentConfigName, accessWithinCluster)
		if err != nil {
			return false, nil
		}

		readyReplicas := deploymentConfig.Status.ReadyReplicas
		numberReplicas := deploymentConfig.Status.Replicas

		fmt.Printf("[DEBUG] Waiting for replicas to be scaled down [%d of %d] (%d seconds elapsed)\n", readyReplicas, numberReplicas, int(time.Since(start).Seconds()))

		if readyReplicas == 0 {
			return true, nil
		}
		return false, nil
	})
}

func ScaleUpDeployment(namespace, deploymentConfigName, accessWithinCluster string, size int32, t int) error {
	client, err := getDeploymentClient(accessWithinCluster)
	if err != nil {
		return err
	}

	deployment, err := GetDeployment(namespace, deploymentConfigName, accessWithinCluster)
	if err != nil {
		return err
	}

	deployment.Spec.Replicas = &size

	_, err = client.Deployments(namespace).Update(deployment)
	if err != nil {
		return err
	}

	var poll = 5 * time.Second
	timeout := time.Duration(t) * time.Second
	start := time.Now()

	fmt.Printf("[DEBUG] Waiting up to %v for deployment to be scaled to %d\n", timeout, size)

	return wait.PollImmediate(poll, timeout, func() (bool, error) {
		deploymentConfig, err := GetDeployment(namespace, deploymentConfigName, accessWithinCluster)
		if err != nil {
			return false, nil
		}

		readyReplicas := deploymentConfig.Status.ReadyReplicas
		fmt.Printf("[DEBUG] Waiting for replicas to be scaled up [%d of %d] (%d seconds elapsed)\n", readyReplicas, size, int(time.Since(start).Seconds()))

		if readyReplicas == size {
			return true, nil
		}
		return false, nil
	})
}
