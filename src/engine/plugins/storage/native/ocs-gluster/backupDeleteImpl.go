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
	"fossul/src/engine/plugins/pluginUtil"
	"fossul/src/engine/util"
	"strings"
)

func (s storagePlugin) BackupDelete(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

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

	backups, err := pluginUtil.ListSnapshots(snapshotList, config.StoragePluginParameters["PvcName"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)

		return result
	}

	backupsByPolicy := util.GetBackupsByPolicy(config.SelectedBackupPolicy, backups)
	backupCount := len(backupsByPolicy)

	if backupCount > config.SelectedBackupRetention {
		count := 1
		for backup := range pluginUtil.ReverseBackupList(backupsByPolicy) {
			if count > config.SelectedBackupRetention {
				msg := util.SetMessage("INFO", fmt.Sprintf("Number of backups [%d] greater than backup retention [%d]", backupCount, config.SelectedBackupRetention))
				messages = append(messages, msg)
				backupCount = backupCount - 1

				backupName := backup.Name + "_" + backup.Policy + "_" + backup.WorkflowId + "_" + util.IntToString(backup.Epoch)
				msg = util.SetMessage("INFO", "Deleting backup "+backupName)
				messages = append(messages, msg)

				var deleteSnapshot []string
				deleteSnapshot = append(deleteSnapshot, "/usr/sbin/gluster")
				deleteSnapshot = append(deleteSnapshot, "--mode=script")
				deleteSnapshot = append(deleteSnapshot, "snapshot")
				deleteSnapshot = append(deleteSnapshot, "delete")
				deleteSnapshot = append(deleteSnapshot, backupName)

				deleteSnapshotResult := k8s.ExecuteCommand(podName, config.StoragePluginParameters["ContainerName"], config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["AccessWithinCluster"], deleteSnapshot...)
				if deleteSnapshotResult.Code != 0 {
					messages = util.PrependMessages(messages, deleteSnapshotResult.Messages)
					result = util.SetResult(1, messages)
					return result
				} else {
					messages = util.PrependMessages(messages, deleteSnapshotResult.Messages)
				}

				msg = util.SetMessage("INFO", "Backup "+backupName+" deleted successfully")
				messages = append(messages, msg)
			}
			count = count + 1
		}
	} else {
		msg := util.SetMessage("INFO", fmt.Sprintf("Backup deletion skipped, there are [%d] backups but backup retention is [%d]", backupCount, config.SelectedBackupRetention))
		messages = append(messages, msg)
	}

	result = util.SetResult(resultCode, messages)
	return result
}
