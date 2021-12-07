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

	snapclient "github.com/fossul/fossul/src/client/k8s/snapshotctrl/client/versioned"
	virtclient "github.com/fossul/fossul/src/client/k8s/virtctrl/client/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	virtv1 "kubevirt.io/api/core/v1"
	snapshotv1 "kubevirt.io/api/snapshot/v1alpha1"
)

func listVirtualMachines(namespace, accessWithinCluster string) (*virtv1.VirtualMachineInstanceList, error) {
	var vms *virtv1.VirtualMachineInstanceList
	virtualMachineClient, err := getVirtualMachineClient(accessWithinCluster)
	if err != nil {
		return vms, err
	}

	vms, err = virtualMachineClient.KubevirtV1().VirtualMachineInstances(namespace).List(context.Background(), metav1.ListOptions{})

	if err != nil {
		return vms, err
	}

	return vms, nil
}

func getVirtualMachines(namespace, accessWithinCluster, vmName string) (*virtv1.VirtualMachine, error) {
	var vm *virtv1.VirtualMachine
	virtualMachineClient, err := getVirtualMachineClient(accessWithinCluster)
	if err != nil {
		return vm, err
	}

	vm, err = virtualMachineClient.KubevirtV1().VirtualMachines(namespace).Get(context.Background(), vmName, metav1.GetOptions{})

	if err != nil {
		return vm, err
	}

	return vm, nil
}

func vmSnapshotSucceeded(vmSnapshot *snapshotv1.VirtualMachineSnapshot) bool {
	return vmSnapshot.Status != nil && vmSnapshot.Status.Phase == snapshotv1.Succeeded
}

func createSNapshot(namespace, accessWithinCluster, vmName string) *snapshotv1.VirtualMachineSnapshot {
	groupName := "kubevirt.io"

	newSnapshot := func() *snapshotv1.VirtualMachineSnapshot {
		return &snapshotv1.VirtualMachineSnapshot{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "snapshot-" + vmName,
				Namespace: namespace,
			},
			Spec: snapshotv1.VirtualMachineSnapshotSpec{
				Source: corev1.TypedLocalObjectReference{
					APIGroup: &groupName,
					Kind:     "VirtualMachine",
					Name:     vmName,
				},
			},
		}
	}

	return newSnapshot()
}

func getVirtualMachineClient(accessWithinCluster string) (*virtclient.Clientset, error) {
	var virtualMachineClient *virtclient.Clientset

	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return virtualMachineClient, err
	}

	virtualMachineClient, err = virtclient.NewForConfig(kubeConfig)
	if err != nil {
		return virtualMachineClient, err
	}

	return virtualMachineClient, nil
}

func getVirtSnapshotClient(accessWithinCluster string) (*snapclient.Clientset, error) {
	var snapshotClient *snapclient.Clientset

	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return snapshotClient, err
	}

	snapshotClient, err = snapclient.NewForConfig(kubeConfig)
	if err != nil {
		return snapshotClient, err
	}

	return snapshotClient, nil
}
