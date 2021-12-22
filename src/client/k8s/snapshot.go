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
	"time"

	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v4/apis/volumesnapshot/v1"
	snapClient "github.com/kubernetes-csi/external-snapshotter/client/v4/clientset/versioned/typed/volumesnapshot/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func CreateSnapshot(snapshotName, namespace, snapshotClassName, pvcName, accessWithinCluster string, t int) error {

	sclient, err := getSnapshotClient(accessWithinCluster)
	if err != nil {
		return err
	}

	snap := generateSnapshot(snapshotName, namespace, snapshotClassName, pvcName, accessWithinCluster)

	_, err = sclient.VolumeSnapshots(snap.Namespace).Create(context.Background(), snap, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("snapshot with name %v created in %v namespace\n", snap.Name, snap.Namespace)

	var poll = 2 * time.Second
	timeout := time.Duration(t) * time.Second
	name := snap.Name
	start := time.Now()
	fmt.Printf("Waiting up to %v seconds to be in Ready state\n", timeout)

	time.Sleep(poll)
	return wait.PollImmediate(poll, timeout, func() (bool, error) {
		fmt.Printf("waiting for snapshot %s (%d seconds elapsed)\n", snap.Name, int(time.Since(start).Seconds()))
		snaps, err := sclient.VolumeSnapshots(snap.Namespace).Get(context.Background(), name, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Error getting snapshot in namespace: '%s': %v\n", snap.Namespace, err)
			return false, err
		}

		if snaps == nil {
			fmt.Printf("Error getting snapshot status: '%s': %v\n", snap.Namespace, err)
			return false, err
		}

		if *snaps.Status.ReadyToUse {
			return true, nil
		}
		return false, nil
	})
}

func ListSnapshots(namespace, accessWithinCluster string) (*snapshotv1.VolumeSnapshotList, error) {
	var snapshots *snapshotv1.VolumeSnapshotList
	sclient, err := getSnapshotClient(accessWithinCluster)
	if err != nil {
		return snapshots, err
	}

	snapshots, err = sclient.VolumeSnapshots(namespace).List(context.Background(), metav1.ListOptions{})

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

	err = sclient.VolumeSnapshots(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func GetSnapshotHandle(contentName *string, namespace string, accessWithinCluster string) (string, error) {
	var snapshotHandle string

	sclient, err := getSnapshotClient(accessWithinCluster)
	if err != nil {
		return snapshotHandle, err
	}

	snapshotContent, err := sclient.VolumeSnapshotContents().Get(context.Background(), *contentName, metav1.GetOptions{})
	if err != nil {
		return snapshotHandle, err
	}

	snapshotHandle = *snapshotContent.Status.SnapshotHandle

	return snapshotHandle, nil
}

func GetSnapshot(name, namespace, accessWithinCluster string) (*snapshotv1.VolumeSnapshot, error) {
	var snapshot *snapshotv1.VolumeSnapshot

	sclient, err := getSnapshotClient(accessWithinCluster)
	if err != nil {
		return snapshot, err
	}

	snapshot, err = sclient.VolumeSnapshots(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return snapshot, err
	}

	return snapshot, nil
}

func GetVolumeSnapshotClassName(storageDriverName, accessWithinCluster string) (string, error) {
	sclient, err := getSnapshotClient(accessWithinCluster)
	if err != nil {
		return "", err
	}

	snapshotClassList, err := sclient.VolumeSnapshotClasses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	for _, snapshotClass := range snapshotClassList.Items {
		if snapshotClass.Driver == storageDriverName {
			return snapshotClass.Name, nil
		}
	}

	return "", errors.New("ERROR: COuldn't determine snapshot class using driver [" + storageDriverName + "]")
}

func getSnapshotClient(accessWithinCluster string) (*snapClient.SnapshotV1Client, error) {
	var sclient *snapClient.SnapshotV1Client

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

func generateSnapshot(snapshotName, namespace, snapshotClassName, pvcName, accessWithinCluster string) *snapshotv1.VolumeSnapshot {

	claimName, _ := GetPersistentVolumeClaim(namespace, pvcName, accessWithinCluster)

	snapshot := snapshotv1.VolumeSnapshot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      snapshotName,
			Namespace: namespace,
		},
		Spec: snapshotv1.VolumeSnapshotSpec{
			VolumeSnapshotClassName: &snapshotClassName,
			Source: snapshotv1.VolumeSnapshotSource{
				PersistentVolumeClaimName: &claimName.Name,
			},
		},
	}

	return &snapshot
}
