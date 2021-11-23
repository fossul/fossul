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
package main

import (
	"github.com/fossul/fossul/src/client/k8s"
	"github.com/fossul/fossul/src/engine/util"
	"github.com/fossul/fossul/src/plugins/pluginUtil"
)

func (s storagePlugin) Restore(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "Performing CSI Ceph snapshot restore")
	messages = append(messages, msg)

	podName, err := k8s.GetPodByName(config.StoragePluginParameters["CephStorageNamespace"], config.StoragePluginParameters["CephToolsPodName"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	snapshots, err := k8s.ListSnapshots(config.StoragePluginParameters["Namespace"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)

		return result
	}

	var snapshotList []string
	for _, snapshot := range snapshots.Items {
		snapshotList = append(snapshotList, snapshot.Name)
	}

	restoreSnapshot := util.GetRestoreSnapshot(config, snapshotList)

	pvName, err := k8s.GetPersistentVolumeName(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["PvcName"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	pv, err := k8s.GetPersistentVolume(pvName, config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	cephVolumeHandle := pv.Spec.CSI.VolumeHandle
	cephVolumeName := pluginUtil.CephVolumeName(cephVolumeHandle)

	snapshot, err := k8s.GetSnapshot(restoreSnapshot, config.StoragePluginParameters["Namespace"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	cephSnapshotHandle, err := k8s.GetSnapshotHandle(snapshot.Status.BoundVolumeSnapshotContentName, config.StoragePluginParameters["Namespace"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	cephSnapshotName := pluginUtil.CephSnapshotName(cephSnapshotHandle)

	var snapshotRestore []string
	snapshotRestore = append(snapshotRestore, "rbd")
	snapshotRestore = append(snapshotRestore, "snap")
	snapshotRestore = append(snapshotRestore, "rollback")
	snapshotRestore = append(snapshotRestore, config.StoragePluginParameters["CephStoragePool"]+"/"+cephVolumeName+"@"+cephSnapshotName)

	snapshotRestoreResult := k8s.ExecuteCommand(podName, config.StoragePluginParameters["CephToolsContainerName"], config.StoragePluginParameters["CephStorageNamespace"], config.AccessWithinCluster, snapshotRestore...)
	if snapshotRestoreResult.Code != 0 {
		messages = util.PrependMessages(messages, snapshotRestoreResult.Messages)
		result = util.SetResult(1, messages)
	} else {
		messages = util.PrependMessages(messages, snapshotRestoreResult.Messages)
	}

	result = util.SetResult(resultCode, messages)
	return result
}
