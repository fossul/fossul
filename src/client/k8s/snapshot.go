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
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
	"github.com/kubernetes-csi/external-snapshotter/pkg/apis/volumesnapshot/v1alpha1"
	snapClient "github.com/kubernetes-csi/external-snapshotter/pkg/client/clientset/versioned/typed/volumesnapshot/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var poll = 2 * time.Second

func createSnapshot(snap *v1alpha1.VolumeSnapshot, t int, accessWithinCluster string) error {

	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return err
	}

	sclient, err := snapClient.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}
	
	_, err = sclient.VolumeSnapshots(snap.Namespace).Create(snap)
	if err != nil {
		return err
	}
	fmt.Printf("snapshot with name %v created in %v namespace", snap.Name, snap.Namespace)

	timeout := time.Duration(t) * time.Minute
	name := snap.Name
	start := time.Now()
	fmt.Printf("Waiting up to %v to be in Ready state", snap)

	return wait.PollImmediate(poll, timeout, func() (bool, error) {
		fmt.Printf("waiting for snapshot %s (%d seconds elapsed)", snap.Name, int(time.Since(start).Seconds()))
		snaps, err := sclient.VolumeSnapshots(snap.Namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Error getting snapshot in namespace: '%s': %v", snap.Namespace, err)
			return false, err
		}
		if snaps.Status.ReadyToUse {
			return true, nil
		}
		return false, nil
	})
}
