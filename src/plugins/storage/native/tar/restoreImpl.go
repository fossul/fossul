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

func (s storagePlugin) Restore(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message

	kubeCmd := GetKubeCmd(config.StoragePluginParameters["KubeCmd"])

	var restorePath string
	if config.StoragePluginParameters["RestorePath"] == "" {
		restorePath = "/"
	} else {
		restorePath = config.StoragePluginParameters["RestorePath"]
	}

	backupPath := "data/backups/" + config.ProfileName + "/" + config.ConfigName
	backupName, err := GetBackupNameFromWorkflowId(util.IntToString(config.SelectedWorkflowId), backupPath)

	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}
	destFile := backupPath + "/" + backupName

	msg := util.SetMessage("INFO", "Performing restore using tar")
	messages = append(messages, msg)

	podName, err := k8s.GetPodName(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["PodSelector"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	msg = util.SetMessage("INFO", "Performing restore for pod "+podName)
	messages = append(messages, msg)

	copyCmd := GetRestoreCopyCommand(config.StoragePluginParameters["Namespace"], podName, kubeCmd, destFile)

	fmt.Println(copyCmd)
	copyCmdResult := util.ExecuteCommand(copyCmd...)
	if copyCmdResult.Code != 0 {
		messages = util.PrependMessages(messages, copyCmdResult.Messages)
		result = util.SetResult(1, messages)
		return result
	} else {
		messages = util.PrependMessages(messages, copyCmdResult.Messages)
	}

	untarCmd := GetRestoreTarCommand(config.StoragePluginParameters["Namespace"], podName, kubeCmd, "/tmp/"+backupName, restorePath)

	untarCmdResult := util.ExecuteCommand(untarCmd...)
	if untarCmdResult.Code != 0 {
		messages = util.PrependMessages(messages, untarCmdResult.Messages)
		result = util.SetResult(1, messages)
		return result
	} else {
		messages = util.PrependMessages(messages, untarCmdResult.Messages)
	}

	cleanupCmd := GetTarCleanupCommand(config.StoragePluginParameters["Namespace"], podName, kubeCmd, backupName)

	cleanupCmdResult := util.ExecuteCommand(cleanupCmd...)
	if cleanupCmdResult.Code != 0 {
		messages = util.PrependMessages(messages, cleanupCmdResult.Messages)
		result = util.SetResult(1, messages)
		return result
	} else {
		messages = util.PrependMessages(messages, cleanupCmdResult.Messages)
	}

	result = util.SetResult(0, messages)
	return result
}
