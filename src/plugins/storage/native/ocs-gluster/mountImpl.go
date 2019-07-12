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
	"fossul/src/client/k8s"
	"fossul/src/plugins/pluginUtil"
	"fossul/src/engine/util"
)

func (s storagePlugin) Mount(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	timestampToString := fmt.Sprintf("%d", config.WorkflowTimestamp)
	backupName := util.GetBackupName(config.StoragePluginParameters["PvcName"], config.SelectedBackupPolicy, config.WorkflowId, timestampToString)

	msg := util.SetMessage("INFO", "Mounting snapshot ["+backupName+"]")
	messages = append(messages, msg)

	podName, err := k8s.GetPodByName(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["PodName"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	podIp, err := k8s.GetPodIp(config.StoragePluginParameters["Namespace"], podName, config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	pvName, err := k8s.GetPersistentVolumeName(config.StoragePluginParameters["DatabaseNamespace"], config.StoragePluginParameters["PvcName"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	glusterVolume, err := k8s.GetGlusterPersistentVolumePath(pvName, config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	var createDir []string
	createDir = append(createDir, "/usr/bin/mkdir")
	createDir = append(createDir, "-p")
	createDir = append(createDir, "/tmp/"+backupName)

	createDirResult := k8s.ExecuteCommand(podName, config.StoragePluginParameters["ContainerName"], config.StoragePluginParameters["Namespace"], config.AccessWithinCluster, createDir...)
	if createDirResult.Code != 0 {
		messages = util.PrependMessages(messages, createDirResult.Messages)
		result = util.SetResult(1, messages)
		return result
	} else {
		messages = util.PrependMessages(messages, createDirResult.Messages)
	}

	var mountSnapshot []string
	mountSnapshot = append(mountSnapshot, "/usr/bin/mount")
	mountSnapshot = append(mountSnapshot, "-t")
	mountSnapshot = append(mountSnapshot, "glusterfs")
	mountSnapshot = append(mountSnapshot, podIp+":/snaps/"+backupName+"/"+glusterVolume)
	mountSnapshot = append(mountSnapshot, "/tmp/"+backupName)

	mountSnapshotResult := k8s.ExecuteCommand(podName, config.StoragePluginParameters["ContainerName"], config.StoragePluginParameters["Namespace"], config.AccessWithinCluster, mountSnapshot...)
	if mountSnapshotResult.Code != 0 {
		messages = util.PrependMessages(messages, mountSnapshotResult.Messages)
		result = util.SetResult(1, messages)

		var deleteDir []string
		deleteDir = append(deleteDir, "/usr/bin/rmdir")
		deleteDir = append(deleteDir, "/tmp/"+backupName)

		deleteDirResult := k8s.ExecuteCommand(podName, config.StoragePluginParameters["ContainerName"], config.StoragePluginParameters["Namespace"], config.AccessWithinCluster, deleteDir...)
		if deleteDirResult.Code != 0 {
			messages = util.PrependMessages(messages, deleteDirResult.Messages)
			result = util.SetResult(1, messages)
			return result
		} else {
			messages = util.PrependMessages(messages, deleteDirResult.Messages)
		}

		return result
	} else {
		messages = util.PrependMessages(messages, mountSnapshotResult.Messages)
	}

	backupPath := config.StoragePluginParameters["BackupDestPath"] + "/" + config.ProfileName + "/" + config.ConfigName
	msg = util.SetMessage("INFO", "Backup name is "+backupName+", Backup path is "+backupPath)
	messages = append(messages, msg)

	err = pluginUtil.CreateDir(backupPath, 0755)
	if err != nil {
		msg = util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)
		return result
	}

	var backupSrcFilePath string
	if config.StoragePluginParameters["SnapshotSubDir"] != "" {
		backupSrcFilePath = "/tmp/" + backupName + "/" + config.StoragePluginParameters["SnapshotSubDir"]

		err = pluginUtil.CreateDir(backupPath+"/"+backupName, 0755)
		if err != nil {
			msg = util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)
			result = util.SetResult(1, messages)
			return result
		}
	} else {
		backupSrcFilePath = "/tmp/" + backupName
	}

	var args []string
	args = append(args, config.StoragePluginParameters["CopyCmdPath"])
	if config.ContainerPlatform == "openshift" {
		args = append(args, "rsync")
		args = append(args, "-n")
		args = append(args, config.StoragePluginParameters["Namespace"])
		args = append(args, podName+":"+backupSrcFilePath)
	} else if config.ContainerPlatform == "kubernetes" {
		args = append(args, "cp")
		args = append(args, config.StoragePluginParameters["Namespace"]+"/"+podName+":"+config.StoragePluginParameters["BackupSrcPath"])
	} else {
		msg = util.SetMessage("ERROR", "Incorrect parameter set for ContainerPlatform ["+config.ContainerPlatform+"]")
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	if config.StoragePluginParameters["SnapshotSubDir"] != "" {
		args = append(args, backupPath+"/"+backupName)
	} else {
		args = append(args, backupPath)
	}

	cmdResult := util.ExecuteCommand(args...)
	if cmdResult.Code != 0 {
		messages = util.PrependMessages(messages, cmdResult.Messages)
		result = util.SetResult(1, messages)
		return result
	} else {
		messages = util.PrependMessages(cmdResult.Messages, messages)
	}

	result = util.SetResult(resultCode, messages)
	return result

}
