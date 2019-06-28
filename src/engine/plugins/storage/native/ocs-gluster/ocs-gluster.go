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
	"fossul/src/engine/client/k8s"
	"fossul/src/engine/util"
	"github.com/heketi/heketi/client/api/go-client"
	//"strings"
)

type storagePlugin string

var StoragePlugin storagePlugin

func (s storagePlugin) SetEnv(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	result = util.SetResult(resultCode, messages)

	return result
}

func (s storagePlugin) Backup(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	timestampToString := fmt.Sprintf("%d", config.WorkflowTimestamp)
	backupName := util.GetBackupName(config.StoragePluginParameters["BackupName"], config.SelectedBackupPolicy, config.WorkflowId, timestampToString)

	msg := util.SetMessage("INFO", "Performing Gluster snapshot")
	messages = append(messages, msg)

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

	pvName, err := k8s.GetPersistentVolumeName(config.StoragePluginParameters["DatabaseNamespace"], config.StoragePluginParameters["PvcName"], config.StoragePluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	glusterVolume, err := k8s.GetGlusterPersistentVolumePath(pvName, config.StoragePluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	podName, err := k8s.GetPodByName(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["PodName"], config.StoragePluginParameters["AccessWithinCluster"])
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

	createSnapshotResult := k8s.ExecuteCommand(podName, config.StoragePluginParameters["ContainerName"], config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["AccessWithinCluster"], createSnapshot...)
	if createSnapshotResult.Code != 0 {
		return createSnapshotResult
	} else {
		messages = util.PrependMessages(messages, createSnapshotResult.Messages)
	}

	var activateSnapshot []string
	activateSnapshot = append(activateSnapshot, "/usr/sbin/gluster")
	activateSnapshot = append(activateSnapshot, "--mode=script")
	activateSnapshot = append(activateSnapshot, "snapshot")
	activateSnapshot = append(activateSnapshot, "activate")
	activateSnapshot = append(activateSnapshot, backupName)

	activateSnapshotResult := k8s.ExecuteCommand(podName, config.StoragePluginParameters["ContainerName"], config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["AccessWithinCluster"], activateSnapshot...)
	if activateSnapshotResult.Code != 0 {
		return activateSnapshotResult
	} else {
		messages = util.PrependMessages(messages, activateSnapshotResult.Messages)
	}

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) Restore(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Restore ***")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) BackupDelete(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	timestampToString := fmt.Sprintf("%d", config.WorkflowTimestamp)
	backupName := util.GetBackupName(config.StoragePluginParameters["BackupName"], config.SelectedBackupPolicy, config.WorkflowId, timestampToString)

	msg := util.SetMessage("INFO", "*** Backup Delete ***")
	messages = append(messages, msg)

	podName, err := k8s.GetPodByName(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["PodName"], config.StoragePluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	var deleteSnapshot []string
	deleteSnapshot = append(deleteSnapshot, "/usr/sbin/gluster")
	deleteSnapshot = append(deleteSnapshot, "--mode=script")
	deleteSnapshot = append(deleteSnapshot, "snapshot")
	deleteSnapshot = append(deleteSnapshot, "delete")
	deleteSnapshot = append(deleteSnapshot, backupName)

	deleteSnapshotResult := k8s.ExecuteCommand(podName, config.StoragePluginParameters["ContainerName"], config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["AccessWithinCluster"], deleteSnapshot...)
	if deleteSnapshotResult.Code != 0 {
		return deleteSnapshotResult
	} else {
		messages = util.PrependMessages(messages, deleteSnapshotResult.Messages)
	}

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) BackupList(config util.Config) util.Backups {
	var backups util.Backups

	return backups
}

func (s storagePlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "ocs-gluster"
	plugin.Description = "OpenShift Container Storage (Gluster)"
	plugin.Version = "1.0.0"
	plugin.Type = "storage"

	var capabilities []util.Capability
	var backupCap util.Capability
	backupCap.Name = "backup"

	var backupListCap util.Capability
	backupListCap.Name = "backupList"

	var backupDeleteCap util.Capability
	backupDeleteCap.Name = "backupDelete"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities, backupCap, backupListCap, backupDeleteCap, infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func main() {}
