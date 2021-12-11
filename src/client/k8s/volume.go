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
	"fmt"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func GetPersistentVolume(pvName, accessWithinCluster string) (*v1.PersistentVolume, error) {
	var pv *v1.PersistentVolume

	client, err := getClient(accessWithinCluster)
	if err != nil {
		return pv, err
	}

	pv, err = client.CoreV1().PersistentVolumes().Get(context.Background(), pvName, metav1.GetOptions{})
	if err != nil {
		return pv, err
	}

	return pv, nil
}

func GetPersistentVolumeClaim(namespace, pvcName, accessWithinCluster string) (*v1.PersistentVolumeClaim, error) {
	var pvc *v1.PersistentVolumeClaim

	client, err := getClient(accessWithinCluster)
	if err != nil {
		return pvc, err
	}

	pvc, err = client.CoreV1().PersistentVolumeClaims(namespace).Get(context.Background(), pvcName, metav1.GetOptions{})
	if err != nil {
		return pvc, err
	}

	return pvc, nil
}

func CreatePersistentVolumeClaimFromSnapshot(pvcName, pvcSize, snapshotName, namespace, storageClassName, accessWithinCluster string) error {
	client, err := getClient(accessWithinCluster)
	if err != nil {
		return err
	}

	pvc := generatePersistentVolumeClaimFromSnapshot(pvcName, pvcSize, snapshotName, namespace, storageClassName)

	_, err = client.CoreV1().PersistentVolumeClaims(namespace).Create(context.Background(), pvc, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func GetPersistentVolumeName(namespace, pvcName, accessWithinCluster string) (string, error) {
	client, err := getClient(accessWithinCluster)
	if err != nil {
		return "", err
	}

	pvc, err := client.CoreV1().PersistentVolumeClaims(namespace).Get(context.Background(), pvcName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return pvc.Spec.VolumeName, nil
}

func GetGlusterVolumePath(pvName, accessWithinCluster string) (string, error) {
	client, err := getClient(accessWithinCluster)
	if err != nil {
		return "", err
	}

	pv, err := client.CoreV1().PersistentVolumes().Get(context.Background(), pvName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return pv.Spec.Glusterfs.Path, nil
}

func GeneratePersistentVolumeClaimVolumeName(pvcName string) *v1.PersistentVolumeClaimVolumeSource {
	volumeSource := v1.PersistentVolumeClaimVolumeSource{
		ClaimName: pvcName,
	}

	return &volumeSource
}

func ListPersistentVolumeClaims(namespace, accessWithinCluster string) (*v1.PersistentVolumeClaimList, error) {
	var pvcList *v1.PersistentVolumeClaimList

	client, err := getClient(accessWithinCluster)
	if err != nil {
		return pvcList, err
	}

	pvcList, err = client.CoreV1().PersistentVolumeClaims(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return pvcList, err
	}

	return pvcList, nil
}

func DeletePersistentVolumeClaim(pvcName, namespace, accessWithinCluster string) error {
	client, err := getClient(accessWithinCluster)
	if err != nil {
		return err
	}

	err = client.CoreV1().PersistentVolumeClaims(namespace).Delete(context.Background(), pvcName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	var poll = 5 * time.Second
	timeout := time.Duration(60) * time.Second
	start := time.Now()

	fmt.Printf("[DEBUG] Waiting up to %v for pvc [%s] deletion\n", timeout, pvcName)

	return wait.PollImmediate(poll, timeout, func() (bool, error) {
		pvcList, err := ListPersistentVolumeClaims(namespace, accessWithinCluster)
		if err != nil {
			return false, nil
		}

		for _, pvc := range pvcList.Items {
			if strings.Compare(pvc.Name, pvcName) == 0 {
				fmt.Printf("[DEBUG] PVC [%s] exists, waiting to be deleted (%d seconds elapsed)\n", pvcName, int(time.Since(start).Seconds()))
				return false, nil
			}
		}

		return true, nil
	})
}

func generatePersistentVolumeClaimFromSnapshot(pvcName, pvcSize, snapshotName, namespace, storageClassName string) *v1.PersistentVolumeClaim {
	apiGroup := "snapshot.storage.k8s.io"
	pvc := v1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcName,
			Namespace: namespace,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			DataSource: &v1.TypedLocalObjectReference{
				Name:     snapshotName,
				Kind:     "VolumeSnapshot",
				APIGroup: &apiGroup,
			},
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(pvcSize),
				},
			},
		},
	}

	return &pvc
}
