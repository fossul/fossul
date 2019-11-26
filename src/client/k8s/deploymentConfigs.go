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
	"fmt"
	apps "github.com/openshift/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

func GetDeploymentConfig(namespace, deploymentConfigName, accessWithinCluster string) (*apps.DeploymentConfig, error) {
	var deploymentConfig *apps.DeploymentConfig

	client, err := getDeploymentConfigClient(accessWithinCluster)
	if err != nil {
		return deploymentConfig, err
	}

	deploymentConfig, err = client.DeploymentConfigs(namespace).Get(deploymentConfigName, metav1.GetOptions{})
	if err != nil {
		return deploymentConfig, err
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

func ScaleDownDeploymentConfig(namespace, deploymentConfigName, accessWithinCluster string, size int32, t int) error {
	client, err := getDeploymentConfigClient(accessWithinCluster)
	if err != nil {
		return err
	}

	deploymentConfig, err := GetDeploymentConfig(namespace, deploymentConfigName, accessWithinCluster)
	if err != nil {
		return err
	}

	deploymentConfig.Spec.Replicas = size

	_, err = client.DeploymentConfigs(namespace).Update(deploymentConfig)
	if err != nil {
		return err
	}

	var poll = 5 * time.Second
	timeout := time.Duration(t) * time.Second
	start := time.Now()

	fmt.Printf("[DEBUG] Waiting up to %v for deployment to be scaled to %d\n", timeout, size)

	return wait.PollImmediate(poll, timeout, func() (bool, error) {
		deploymentConfig, err := GetDeploymentConfig(namespace, deploymentConfigName, accessWithinCluster)
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

func ScaleUpDeploymentConfig(namespace, deploymentConfigName, accessWithinCluster string, size int32, t int) error {
	client, err := getDeploymentConfigClient(accessWithinCluster)
	if err != nil {
		return err
	}

	deploymentConfig, err := GetDeploymentConfig(namespace, deploymentConfigName, accessWithinCluster)
	if err != nil {
		return err
	}

	deploymentConfig.Spec.Replicas = size

	_, err = client.DeploymentConfigs(namespace).Update(deploymentConfig)
	if err != nil {
		return err
	}

	var poll = 5 * time.Second
	timeout := time.Duration(t) * time.Second
	start := time.Now()

	fmt.Printf("[DEBUG] Waiting up to %v for deployment to be scaled to %d\n", timeout, size)

	return wait.PollImmediate(poll, timeout, func() (bool, error) {
		deploymentConfig, err := GetDeploymentConfig(namespace, deploymentConfigName, accessWithinCluster)
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

func UpdateDeploymentConfigVolume(pvcName, namespace, deploymentConfigName, accessWithinCluster string) error {

	client, err := getDeploymentConfigClient(accessWithinCluster)
	if err != nil {
		return err
	}

	deploymentConfig, err := GetDeploymentConfig(namespace, deploymentConfigName, accessWithinCluster)
	if err != nil {
		return err
	}

	//volumes := deploymentConfig.Spec.Template.Spec.Volumes
	//for _, volume := range volumes {
	// fmt.Println("[DEBUG] Updating pv [" + volume.Name + "] with new pvc [" + pvcName + "]")
	//	volume.PersistentVolumeClaim = GeneratePersistentVolumeClaimVolumeName(pvcName)
	//}

	deploymentConfig.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim = GeneratePersistentVolumeClaimVolumeName(pvcName)

	_, err = client.DeploymentConfigs(namespace).Update(deploymentConfig)
	if err != nil {
		return err
	}

	return nil
}
