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
	"strings"

	"github.com/fossul/fossul/src/client/k8s"
	"github.com/fossul/fossul/src/engine/util"
	"github.com/fossul/fossul/src/plugins/pluginUtil"
)

func (s storagePlugin) Backup(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message

	var backup util.Backup

	kubeCmd := GetKubeCmd(config.StoragePluginParameters["KubeCmd"])

	backup.Epoch = int(config.WorkflowTimestamp)
	backup.Policy = config.SelectedBackupPolicy
	backup.WorkflowId = config.WorkflowId

	timestampToString := fmt.Sprintf("%d", config.WorkflowTimestamp)
	backup.Timestamp = timestampToString

	backupName := util.GetBackupName(config.StoragePluginParameters["BackupName"], config.SelectedBackupPolicy, config.WorkflowId, timestampToString)
	backupPath := "data/backups/" + config.ProfileName + "/" + config.ConfigName
	destPath := backupPath + "/" + backupName + ".tar.gz"
	backupSrcFilePaths := getBackupSrcPaths(config)

	msg := util.SetMessage("INFO", "Performing container backup using tar")
	messages = append(messages, msg)

	podName, err := k8s.GetPodName(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["PodSelector"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	msg = util.SetMessage("INFO", "Performing backup for pod "+podName)
	messages = append(messages, msg)

	msg = util.SetMessage("INFO", "Backup name is "+backupName+", Backup path is "+backupPath)
	messages = append(messages, msg)

	err = pluginUtil.CreateDir(backupPath, 0755)
	if err != nil {
		msg = util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)
		return result
	}

	contents, tarCmd := GetTarBackupCommand(config.StoragePluginParameters["Namespace"], podName, backupName+".tar.gz", kubeCmd, backupSrcFilePaths)
	cmdResult := util.ExecuteCommand(tarCmd...)
	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	copyCmd := GetBackupCopyCommand(config.StoragePluginParameters["Namespace"], podName, kubeCmd, backupName+".tar.gz", destPath)
	cmdResult = util.ExecuteCommand(copyCmd...)
	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	cleanupCmd := GetTarCleanupCommand(config.StoragePluginParameters["Namespace"], podName, kubeCmd, backupName+".tar.gz")
	cmdResult = util.ExecuteCommand(cleanupCmd...)
	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	result = util.SetResult(0, messages)
	backup.Contents = contents
	result.Backup = backup

	return result
}

func getBackupSrcPaths(config util.Config) []string {
	var backupSrcFilePaths []string

	if config.AutoDiscovery == true {
		dataPaths := strings.Split(config.StoragePluginParameters["DataFilePaths"], ",")
		logPaths := strings.Split(config.StoragePluginParameters["LogFilePaths"], ",")

		for _, dataPath := range dataPaths {
			if dataPath == "" {
				continue
			}

			backupSrcFilePaths = append(backupSrcFilePaths, dataPath)
		}

		for _, logPath := range logPaths {
			if logPath == "" {
				continue
			}

			backupSrcFilePaths = append(backupSrcFilePaths, logPath)
		}
	} else {
		backupSrcFilePaths = strings.Split(config.StoragePluginParameters["BackupSrcPaths"], ",")
	}

	return backupSrcFilePaths
}
