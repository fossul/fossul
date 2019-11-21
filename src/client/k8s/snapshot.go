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
	"github.com/kubernetes-csi/external-snapshotter/pkg/apis/volumesnapshot/v1alpha1"
	snapClient "github.com/kubernetes-csi/external-snapshotter/pkg/client/clientset/versioned/typed/volumesnapshot/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

var poll = 2 * time.Second

func CreateSnapshot(snapshotName, namespace, snapshotClassName, pvcName, accessWithinCluster string, t int) error {

	sclient, err := getSnapshotClient(accessWithinCluster)
	if err != nil {
		return err
	}

	snap := generateSnapshot(snapshotName, namespace, snapshotClassName, pvcName)

	_, err = sclient.VolumeSnapshots(snap.Namespace).Create(snap)
	if err != nil {
		return err
	}

	fmt.Printf("snapshot with name %v created in %v namespace\n", snap.Name, snap.Namespace)

	timeout := time.Duration(t) * time.Second
	name := snap.Name
	start := time.Now()
	fmt.Printf("Waiting up to %v to be in Ready state\n", snap)

	return wait.PollImmediate(poll, timeout, func() (bool, error) {
		fmt.Printf("waiting for snapshot %s (%d seconds elapsed)\n", snap.Name, int(time.Since(start).Seconds()))
		snaps, err := sclient.VolumeSnapshots(snap.Namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Error getting snapshot in namespace: '%s': %v\n", snap.Namespace, err)
			return false, err
		}
		if snaps.Status.ReadyToUse {
			return true, nil
		}
		return false, nil
	})
}

func ListSnapshots(namespace, accessWithinCluster string) (*v1alpha1.VolumeSnapshotList, error) {
	var snapshots *v1alpha1.VolumeSnapshotList
	sclient, err := getSnapshotClient(accessWithinCluster)
	if err != nil {
		return snapshots, err
	}

	snapshots, err = sclient.VolumeSnapshots(namespace).List(metav1.ListOptions{})
	if err != nil {
		return snapshots, err
	}

	return snapshots, nil
}

func DeleteSnapshot(name, namespace, accessWithinCluster string) error {
	sclient, err := getSnapshotClient(accessWithinCluster)
	if err != nil {
		return err
	}

	err = sclient.VolumeSnapshots(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func getSnapshotClient(accessWithinCluster string) (*snapClient.SnapshotV1alpha1Client, error) {
	var sclient *snapClient.SnapshotV1alpha1Client
	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return sclient, err
	}

	sclient, err = snapClient.NewForConfig(kubeConfig)
	if err != nil {
		return sclient, err
	}

	return sclient, nil
}

func generateSnapshot(snapshotName, namespace, snapshotClassName, pvcName string) *v1alpha1.VolumeSnapshot {
	snapshot := v1alpha1.VolumeSnapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      snapshotName,
			Namespace: namespace,
		},
		Spec: v1alpha1.VolumeSnapshotSpec{
			VolumeSnapshotClassName: &snapshotClassName,
			Source: &v1.TypedLocalObjectReference{
				Name: pvcName,
				Kind: "PersistentVolumeClaim",
			},
		},
	}

	return &snapshot
}
