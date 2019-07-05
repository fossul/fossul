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
)

func (s storagePlugin) Unmount(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	timestampToString := fmt.Sprintf("%d", config.WorkflowTimestamp)
	backupName := util.GetBackupName(config.StoragePluginParameters["PvcName"], config.SelectedBackupPolicy, config.WorkflowId, timestampToString)

	msg := util.SetMessage("INFO", "Mounting snapshot ["+backupName+"]")
	messages = append(messages, msg)

	podName, err := k8s.GetPodByName(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["PodName"], config.StoragePluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	var unmountSnapshot []string
	unmountSnapshot = append(unmountSnapshot, "/usr/bin/umount")
	unmountSnapshot = append(unmountSnapshot, "/tmp/"+backupName)

	unmountSnapshotResult := k8s.ExecuteCommand(podName, config.StoragePluginParameters["ContainerName"], config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["AccessWithinCluster"], unmountSnapshot...)
	if unmountSnapshotResult.Code != 0 {
		messages = util.PrependMessages(messages, unmountSnapshotResult.Messages)
		result = util.SetResult(1, messages)
		return result
	} else {
		messages = util.PrependMessages(messages, unmountSnapshotResult.Messages)
	}

	var deleteDir []string
	deleteDir = append(deleteDir, "/usr/bin/rmdir")
	deleteDir = append(deleteDir, "/tmp/"+backupName)

	deleteDirResult := k8s.ExecuteCommand(podName, config.StoragePluginParameters["ContainerName"], config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["AccessWithinCluster"], deleteDir...)
	if deleteDirResult.Code != 0 {
		messages = util.PrependMessages(messages, deleteDirResult.Messages)
		result = util.SetResult(1, messages)
		return result
	} else {
		messages = util.PrependMessages(messages, deleteDirResult.Messages)
	}

	result = util.SetResult(resultCode, messages)
	return result
}
