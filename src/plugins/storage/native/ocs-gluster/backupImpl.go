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
	"fmt"
	"github.com/fossul/fossul/src/client/k8s"
	"github.com/fossul/fossul/src/engine/util"
)

func (s storagePlugin) Backup(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	timestampToString := fmt.Sprintf("%d", config.WorkflowTimestamp)
	backupName := util.GetBackupName(config.StoragePluginParameters["PvcName"], config.SelectedBackupPolicy, config.WorkflowId, timestampToString)

	msg := util.SetMessage("INFO", "Performing Gluster snapshot")
	messages = append(messages, msg)

	/* Heketi example, might need heketi so keeping code for now
	heketi := client.NewClient(config.StoragePluginParameters["HeketiUrl"], config.StoragePluginParameters["HeketiUser"], config.StoragePluginParameters["HeketiToken"])

	list, err := heketi.VolumeList()
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	for _, volumeId := range list.Volumes {
		volumeInfo, err := heketi.VolumeInfo(volumeId)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		fmt.Println("DEBUG:", volumeInfo.Name, volumeInfo.Mount)
	}
	*/

	pvName, err := k8s.GetPersistentVolumeName(config.StoragePluginParameters["DatabaseNamespace"], config.StoragePluginParameters["PvcName"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	glusterVolume, err := k8s.GetGlusterVolumePath(pvName, config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	podName, err := k8s.GetPodByName(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["PodName"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	msg = util.SetMessage("INFO", "Performing backup of pod ["+podName+"] pv ["+pvName+"] gluster volume ["+glusterVolume+"]")
	messages = append(messages, msg)

	var createSnapshot []string
	createSnapshot = append(createSnapshot, "/usr/sbin/gluster")
	createSnapshot = append(createSnapshot, "--mode=script")
	createSnapshot = append(createSnapshot, "snapshot")
	createSnapshot = append(createSnapshot, "create")
	createSnapshot = append(createSnapshot, backupName)
	createSnapshot = append(createSnapshot, glusterVolume)
	createSnapshot = append(createSnapshot, "no-timestamp")

	createSnapshotResult := k8s.ExecuteCommand(podName, config.StoragePluginParameters["ContainerName"], config.StoragePluginParameters["Namespace"], config.AccessWithinCluster, createSnapshot...)
	if createSnapshotResult.Code != 0 {
		messages = util.PrependMessages(messages, createSnapshotResult.Messages)
		result = util.SetResult(1, messages)
		return result
	} else {
		messages = util.PrependMessages(messages, createSnapshotResult.Messages)
	}

	var activateSnapshot []string
	activateSnapshot = append(activateSnapshot, "/usr/sbin/gluster")
	activateSnapshot = append(activateSnapshot, "--mode=script")
	activateSnapshot = append(activateSnapshot, "snapshot")
	activateSnapshot = append(activateSnapshot, "activate")
	activateSnapshot = append(activateSnapshot, backupName)

	activateSnapshotResult := k8s.ExecuteCommand(podName, config.StoragePluginParameters["ContainerName"], config.StoragePluginParameters["Namespace"], config.AccessWithinCluster, activateSnapshot...)
	if activateSnapshotResult.Code != 0 {
		messages = util.PrependMessages(messages, activateSnapshotResult.Messages)
		result = util.SetResult(1, messages)
		return result
	} else {
		messages = util.PrependMessages(messages, activateSnapshotResult.Messages)
	}

	result = util.SetResult(resultCode, messages)
	return result
}
