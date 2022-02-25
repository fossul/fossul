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
	"context"
	"errors"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPod(podName, namespace, accessWithinCluster string) (*v1.Pod, error) {
	var pod *v1.Pod
	client, err := getClient(accessWithinCluster)
	if err != nil {
		return pod, err
	}

	pod, err = client.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})

	if err != nil {
		return pod, err
	}

	return pod, nil
}

func GetPodContainer(podName, namespace, accessWithinCluster string) (string, error) {
	var containerName string
	var pod *v1.Pod
	client, err := getClient(accessWithinCluster)
	if err != nil {
		return containerName, err
	}

	pod, err = client.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if pod.Spec.Containers != nil {
		if len(pod.Spec.Containers) > 1 {
			return "", errors.New("Pod [" + pod.Name + "] has more than one containers please specify ContainerName parameter in storage plugin config")
		}
		containerName = pod.Spec.Containers[0].Name
	} else {
		return "", errors.New("Couldn't get ContainerName from pod [" + pod.Name + "] Spec, try providing ContainerName parameter in storage plugin config")
	}

	return containerName, nil
}

func GetPodName(namespace, podSelector, accessWithinCluster string) (string, error) {
	client, err := getClient(accessWithinCluster)
	if err != nil {
		return "", err
	}

	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	fmt.Println("[INFO]: Pods in namespace", namespace)
	var ourPod string
	for _, pod := range pods.Items {
		fmt.Println("[INFO] Pod", pod.Name, pod.Status.Phase)
		if strings.Contains(pod.Name, podSelector) && pod.Status.Phase == "Running" {
			fmt.Println("[INFO] Running Pod Found:", pod.Name)
			ourPod = pod.Name
		}
	}

	return ourPod, nil
}

func GetPodByName(namespace, podName, accessWithinCluster string) (string, error) {
	client, err := getClient(accessWithinCluster)
	if err != nil {
		return "", err
	}

	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	fmt.Println("[INFO]: Pods in namespace", namespace)
	var ourPod string
	for _, pod := range pods.Items {
		fmt.Println("[INFO] Pod", pod.Name, pod.Status.Phase)
		if strings.Contains(pod.Name, podName) && pod.Status.Phase == "Running" {
			fmt.Println("[INFO] Running Pod Found:", pod.Name)
			ourPod = pod.Name
		}
	}

	return ourPod, nil
}

func GetPodIp(namespace, podName, accessWithinCluster string) (string, error) {
	client, err := getClient(accessWithinCluster)
	if err != nil {
		return "", err
	}

	pod, err := client.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	podIp := pod.Status.PodIP
	return podIp, nil
}
