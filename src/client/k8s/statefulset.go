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
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func GetStatefulSet(namespace, name, accessWithinCluster string) (*apps.StatefulSet, error) {
	var statefulset *apps.StatefulSet
	client, err := getAppsClient(accessWithinCluster)
	if err != nil {
		return statefulset, err
	}

	statefulset, err = client.StatefulSets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return statefulset, nil
}

func GetStatefulSetScaleInteger(namespace, name, accessWithinCluster string) (*int32, error) {
	statefulset, err := GetStatefulSet(namespace, name, accessWithinCluster)
	if err != nil {
		return nil, err
	}

	return statefulset.Spec.Replicas, nil
}

func ScaleDownStatefulSet(namespace, name, accessWithinCluster string, size int32, t int) error {
	client, err := getAppsClient(accessWithinCluster)
	if err != nil {
		return err
	}

	statefulset, err := GetStatefulSet(namespace, name, accessWithinCluster)
	if err != nil {
		return err
	}

	statefulset.Spec.Replicas = &size

	_, err = client.StatefulSets(namespace).Update(context.Background(), statefulset, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	var poll = 5 * time.Second
	timeout := time.Duration(t) * time.Second
	start := time.Now()

	fmt.Printf("[DEBUG] Waiting up to %v for deployment to be scaled to %d\n", timeout, size)

	return wait.PollImmediate(poll, timeout, func() (bool, error) {
		statefulsetConfig, err := GetStatefulSet(namespace, name, accessWithinCluster)
		if err != nil {
			return false, nil
		}

		readyReplicas := statefulsetConfig.Status.ReadyReplicas
		numberReplicas := statefulsetConfig.Status.Replicas

		fmt.Printf("[DEBUG] Waiting for replicas to be scaled down [%d of %d] (%d seconds elapsed)\n", readyReplicas, numberReplicas, int(time.Since(start).Seconds()))

		if readyReplicas == 0 {
			return true, nil
		}
		return false, nil
	})
}

func ScaleUpStatefulSet(namespace, name, accessWithinCluster string, size int32, t int) error {
	client, err := getAppsClient(accessWithinCluster)
	if err != nil {
		return err
	}

	statefulset, err := GetStatefulSet(namespace, name, accessWithinCluster)
	if err != nil {
		return err
	}

	statefulset.Spec.Replicas = &size

	_, err = client.StatefulSets(namespace).Update(context.Background(), statefulset, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	var poll = 5 * time.Second
	timeout := time.Duration(t) * time.Second
	start := time.Now()

	fmt.Printf("[DEBUG] Waiting up to %v for deployment to be scaled to %d\n", timeout, size)

	return wait.PollImmediate(poll, timeout, func() (bool, error) {
		statefulsetConfig, err := GetStatefulSet(namespace, name, accessWithinCluster)
		if err != nil {
			return false, nil
		}

		readyReplicas := statefulsetConfig.Status.ReadyReplicas
		fmt.Printf("[DEBUG] Waiting for replicas to be scaled up [%d of %d] (%d seconds elapsed)\n", readyReplicas, size, int(time.Since(start).Seconds()))

		if readyReplicas == size {
			return true, nil
		}
		return false, nil
	})
}

func UpdateStatefulSetVolume(pvcName, restorePvcName, namespace, deploymentName, accessWithinCluster string) error {

	client, err := getAppsClient(accessWithinCluster)
	if err != nil {
		return err
	}

	statefulset, err := GetStatefulSet(namespace, deploymentName, accessWithinCluster)
	if err != nil {
		return err
	}

	for _, volume := range statefulset.Spec.Template.Spec.Volumes {
		if volume.Name == pvcName {
			volume.PersistentVolumeClaim.ClaimName = restorePvcName
		}
	}

	_, err = client.StatefulSets(namespace).Update(context.Background(), statefulset, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
