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
	"fossul/src/engine/client/k8s"
	"fossul/src/engine/util"
	"strings"
)

func (s storagePlugin) Restore(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "Performing Gluster snapshot restore")
	messages = append(messages, msg)

	podName, err := k8s.GetPodByName(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["PodName"], config.StoragePluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	var listSnapshot []string
	listSnapshot = append(listSnapshot, "/usr/sbin/gluster")
	listSnapshot = append(listSnapshot, "--mode=script")
	listSnapshot = append(listSnapshot, "snapshot")
	listSnapshot = append(listSnapshot, "list")

	listSnapshotResult, listSnapshotStdout := k8s.ExecuteCommandWithStdout(podName, config.StoragePluginParameters["ContainerName"], config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["AccessWithinCluster"], listSnapshot...)
	if listSnapshotResult.Code != 0 {
		messages = util.PrependMessages(messages, listSnapshotResult.Messages)
		result = util.SetResult(1, messages)
		return result
	} else {
		messages = util.PrependMessages(messages, listSnapshotResult.Messages)
	}

	snapshotList := strings.Split(listSnapshotStdout, "\n")
	restoreSnapshot := util.GetRestoreSnapshot(config, snapshotList)

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

	var stopGlusterVolume []string
	stopGlusterVolume = append(stopGlusterVolume, "/usr/sbin/gluster")
	stopGlusterVolume = append(stopGlusterVolume, "--mode=script")
	stopGlusterVolume = append(stopGlusterVolume, "vol")
	stopGlusterVolume = append(stopGlusterVolume, "stop")
	stopGlusterVolume = append(stopGlusterVolume, glusterVolume)

	stopGlusterVolumeResult := k8s.ExecuteCommand(podName, config.StoragePluginParameters["ContainerName"], config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["AccessWithinCluster"], stopGlusterVolume...)
	if stopGlusterVolumeResult.Code != 0 {
		messages = util.PrependMessages(messages, stopGlusterVolumeResult.Messages)
		result = util.SetResult(1, messages)
		return result
	} else {
		messages = util.PrependMessages(messages, stopGlusterVolumeResult.Messages)
	}

	var snapshotRestore []string
	snapshotRestore = append(snapshotRestore, "/usr/sbin/gluster")
	snapshotRestore = append(snapshotRestore, "--mode=script")
	snapshotRestore = append(snapshotRestore, "snapshot")
	snapshotRestore = append(snapshotRestore, "restore")
	snapshotRestore = append(snapshotRestore, restoreSnapshot)

	snapshotRestoreResult := k8s.ExecuteCommand(podName, config.StoragePluginParameters["ContainerName"], config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["AccessWithinCluster"], snapshotRestore...)
	if snapshotRestoreResult.Code != 0 {
		messages = util.PrependMessages(messages, snapshotRestoreResult.Messages)
		result = util.SetResult(1, messages)
		return result
	} else {
		messages = util.PrependMessages(messages, snapshotRestoreResult.Messages)
	}

	var startGlusterVolume []string
	startGlusterVolume = append(startGlusterVolume, "/usr/sbin/gluster")
	startGlusterVolume = append(startGlusterVolume, "--mode=script")
	startGlusterVolume = append(startGlusterVolume, "vol")
	startGlusterVolume = append(startGlusterVolume, "start")
	startGlusterVolume = append(startGlusterVolume, glusterVolume)

	startGlusterVolumeResult := k8s.ExecuteCommand(podName, config.StoragePluginParameters["ContainerName"], config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["AccessWithinCluster"], startGlusterVolume...)
	if startGlusterVolumeResult.Code != 0 {
		messages = util.PrependMessages(messages, startGlusterVolumeResult.Messages)
		result = util.SetResult(1, messages)
		return result
	} else {
		messages = util.PrependMessages(messages, startGlusterVolumeResult.Messages)
	}

	result = util.SetResult(resultCode, messages)
	return result
}
