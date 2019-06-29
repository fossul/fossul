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
	//	"k8s.io/apimachinery/pkg/api/errors"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

func GetPod(namespace, serviceName, accessWithinCluster string) (string, error) {
	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return "", err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return "", err
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	fmt.Println("[INFO]: Pods in namespace", namespace)
	var ourPod string
	for _, pod := range pods.Items {
		fmt.Println("[INFO] Pod", pod.Name, pod.Status.Phase)
		if strings.Contains(pod.Name, serviceName) && pod.Status.Phase == "Running" {
			fmt.Println("[INFO] Running Pod Found:", pod.Name)
			ourPod = pod.Name
		}
	}

	/*
		pod := "fossul-app-2-zpdgr"
		_, err = clientset.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
				pod, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		}
	*/

	return ourPod, nil
}

func GetPodByName(namespace, podName, accessWithinCluster string) (string, error) {
	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return "", err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return "", err
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
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

func GetPersistentVolumeName(namespace, pvcName, accessWithinCluster string) (string, error) {
	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return "", err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return "", err
	}

	pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Get(pvcName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return pvc.Spec.VolumeName, nil
}

func GetGlusterPersistentVolumePath(pvName, accessWithinCluster string) (string, error) {
	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return "", err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return "", err
	}

	pv, err := clientset.CoreV1().PersistentVolumes().Get(pvName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return pv.Spec.Glusterfs.Path, nil
}
