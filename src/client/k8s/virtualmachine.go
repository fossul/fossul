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
	"time"

	snapclient "github.com/fossul/fossul/src/client/k8s/snapshotctrl/client/versioned"
	virtclient "github.com/fossul/fossul/src/client/k8s/virtctrl/client/versioned"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	virtv1 "kubevirt.io/api/core/v1"
	snapshotv1 "kubevirt.io/api/snapshot/v1alpha1"
)

func CreateVirtualMachineSnapshot(namespace, accessWithinCluster, vmName string) error {

	var snapshotClient *snapclient.Clientset
	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return err
	}

	snapshotClient, err = snapclient.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	newSnapshot := getNewSnapshot(namespace, accessWithinCluster, vmName)
	_, err = snapshotClient.SnapshotV1alpha1().VirtualMachineSnapshots(namespace).Create(context.Background(), newSnapshot, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	err = waitForSnapshot(namespace, accessWithinCluster, newSnapshot.Name)
	if err != nil {
		return err
	}

	return nil
}

func DeleteVirtualMachineSnapshot(namespace, accessWithinCluster, snapshotName string) error {
	var snapshotClient *snapclient.Clientset

	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return err
	}

	snapshotClient, err = snapclient.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	err = snapshotClient.SnapshotV1alpha1().VirtualMachineSnapshots(namespace).Delete(context.Background(), snapshotName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

func UpdateVirtualMachineDisk(namespace, accessWithinCluster, vmName, pvcName string) error {
	virtualMachineClient, err := getVirtualMachineClient(accessWithinCluster)
	if err != nil {
		return err
	}

	vm, err := GetVirtualMachine(namespace, accessWithinCluster, vmName)
	if err != nil {
		return err
	}

	// need to persist and map disk to pvc in case we have many
	volumes := vm.Spec.Template.Spec.Volumes
	for _, volume := range volumes {
		if volume.Name == "rootdisk" {
			volume.PersistentVolumeClaim.ClaimName = pvcName
		}
	}

	_, err = virtualMachineClient.KubevirtV1().VirtualMachines(namespace).Update(context.Background(), vm, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func getNewSnapshot(namespace, accessWithinCluster, vmName string) *snapshotv1.VirtualMachineSnapshot {
	groupName := "kubevirt.io"

	newSnapshot := func() *snapshotv1.VirtualMachineSnapshot {
		return &snapshotv1.VirtualMachineSnapshot{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fossul-snapshot-" + vmName,
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

func waitForSnapshot(namespace, accessWithinCluster, snapshotName string) error {
	var snapshotClient *snapclient.Clientset
	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return err
	}

	snapshotClient, err = snapclient.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	var poll = 5 * time.Second
	timeout := time.Duration(300) * time.Second
	start := time.Now()

	fmt.Printf("[DEBUG] Waiting up to %v for vm snapshot to be created\n", timeout)

	return wait.PollImmediate(poll, timeout, func() (bool, error) {
		snapshot, err := snapshotClient.SnapshotV1alpha1().VirtualMachineSnapshots(namespace).Get(context.Background(), snapshotName, metav1.GetOptions{})
		if err != nil {
			return false, nil
		}

		fmt.Printf("[DEBUG] Waiting for vm snapshot [%s] (%d seconds elapsed)\n", snapshotName, int(time.Since(start).Seconds()))

		if *snapshot.Status.ReadyToUse {
			return true, nil
		}

		return false, nil
	})
}

func ListVirtualMachines(namespace, accessWithinCluster string) (*virtv1.VirtualMachineInstanceList, error) {
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

func ListVirtualMachineSnapshots(namespace, accessWithinCluster, vmName string) (*snapshotv1.VirtualMachineSnapshotList, error) {
	var snapshotList *snapshotv1.VirtualMachineSnapshotList
	var snapshotClient *snapclient.Clientset

	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		return snapshotList, err
	}

	snapshotClient, err = snapclient.NewForConfig(kubeConfig)
	if err != nil {
		return snapshotList, err
	}

	snapshotList, err = snapshotClient.SnapshotV1alpha1().VirtualMachineSnapshots(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return snapshotList, err
	}

	return snapshotList, nil

}

func GetVirtualMachine(namespace, accessWithinCluster, vmName string) (*virtv1.VirtualMachine, error) {
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

func StartVirtualMachine(namespace, accessWithinCluster, vmName string) error {
	var vm *virtv1.VirtualMachine
	var vmRunStrategy virtv1.VirtualMachineRunStrategy = "Always"
	virtualMachineClient, err := getVirtualMachineClient(accessWithinCluster)
	if err != nil {
		return err
	}

	vm, err = virtualMachineClient.KubevirtV1().VirtualMachines(namespace).Get(context.Background(), vmName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	//*vm.Spec.Running = true
	vm.Spec.Running = nil
	vm.Spec.RunStrategy = &vmRunStrategy
	_, err = virtualMachineClient.KubevirtV1().VirtualMachines(namespace).Update(context.Background(), vm, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func StopVirtualMachine(namespace, accessWithinCluster, vmName string) error {
	var vm *virtv1.VirtualMachine
	var vmRunStrategy virtv1.VirtualMachineRunStrategy = "Halted"

	virtualMachineClient, err := getVirtualMachineClient(accessWithinCluster)
	if err != nil {
		return err
	}

	vm, err = virtualMachineClient.KubevirtV1().VirtualMachines(namespace).Get(context.Background(), vmName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	//*vm.Spec.Running = false
	vm.Spec.Running = nil
	vm.Spec.RunStrategy = &vmRunStrategy
	_, err = virtualMachineClient.KubevirtV1().VirtualMachines(namespace).Update(context.Background(), vm, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func vmSnapshotSucceeded(vmSnapshot *snapshotv1.VirtualMachineSnapshot) bool {
	return vmSnapshot.Status != nil && vmSnapshot.Status.Phase == snapshotv1.Succeeded
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
