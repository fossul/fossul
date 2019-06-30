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
	backupSrcFilePaths := getBackupSrcPaths(config)

	msg := util.SetMessage("INFO", "Performing container backup")
	messages = append(messages, msg)

	podName, err := k8s.GetPod(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["ServiceName"], config.StoragePluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	msg = util.SetMessage("INFO", "Performing backup for pod "+podName)
	messages = append(messages, msg)

	timestampToString := fmt.Sprintf("%d", config.WorkflowTimestamp)
	backupName := util.GetBackupName(config.StoragePluginParameters["BackupName"], config.SelectedBackupPolicy, config.WorkflowId, timestampToString)
	backupPath := util.GetBackupPathFromConfig(config)
	msg = util.SetMessage("INFO", "Backup name is "+backupName+", Backup path is "+backupPath)
	messages = append(messages, msg)

	err = pluginUtil.CreateDir(backupPath, 0755)
	if err != nil {
		msg = util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)
		return result
	}

	for _, backupSrcFilePath := range backupSrcFilePaths {
		var args []string
		args = append(args, config.StoragePluginParameters["CopyCmdPath"])
		if config.StoragePluginParameters["ContainerPlatform"] == "openshift" {
			args = append(args, "rsync")
			args = append(args, "-n")
			args = append(args, config.StoragePluginParameters["Namespace"])
			args = append(args, podName+":"+backupSrcFilePath)
		} else if config.StoragePluginParameters["ContainerPlatform"] == "kubernetes" {
			args = append(args, "cp")
			args = append(args, config.StoragePluginParameters["Namespace"]+"/"+podName+":"+config.StoragePluginParameters["BackupSrcPath"])
		} else {
			msg = util.SetMessage("ERROR", "Incorrect parameter set for ContainerPlatform ["+config.StoragePluginParameters["ContainerPlatform"]+"]")
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		args = append(args, backupPath)

		cmdResult := util.ExecuteCommand(args...)
		if cmdResult.Code != 0 {
			return cmdResult
		} else {
			messages = util.PrependMessages(cmdResult.Messages, messages)
		}
	}

	result = util.SetResult(0, messages)
	return result
}

func (s storagePlugin) Restore(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message

	msg := util.SetMessage("INFO", "Performing container restore")
	messages = append(messages, msg)

	podName, err := k8s.GetPod(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["ServiceName"], config.StoragePluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	msg = util.SetMessage("INFO", "Performing restore for pod "+podName)
	messages = append(messages, msg)

	restorePath, err := util.GetRestoreSrcPath(config)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)
		return result
	}

	if restorePath == "" {
		msg = util.SetMessage("ERROR", "Restore data no longer available for workflow id ["+util.IntToString(config.SelectedWorkflowId)+"], check retention policy")
		messages = append(messages, msg)
		result = util.SetResult(1, messages)
		return result
	}

	msg = util.SetMessage("INFO", "Restore source path is ["+restorePath+"]")
	messages = append(messages, msg)

	restoreDestPath := "/tmp/" + util.IntToString(config.SelectedWorkflowId)

	var args []string
	args = append(args, config.StoragePluginParameters["CopyCmdPath"])
	if config.StoragePluginParameters["ContainerPlatform"] == "openshift" {
		args = append(args, "rsync")
		args = append(args, "-n")
		args = append(args, config.StoragePluginParameters["Namespace"])
		args = append(args, restorePath)
		args = append(args, podName+":"+restoreDestPath)
	} else if config.StoragePluginParameters["ContainerPlatform"] == "kubernetes" {
		args = append(args, "cp")
		args = append(args, restorePath)
		args = append(args, config.StoragePluginParameters["Namespace"]+"/"+podName+":"+restoreDestPath)
	} else {
		msg = util.SetMessage("ERROR", "Incorrect parameter set for ContainerPlatform ["+config.StoragePluginParameters["ContainerPlatform"]+"]")
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	cmdResult := util.ExecuteCommand(args...)
	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(cmdResult.Messages, messages)
	}

	result = util.SetResult(0, messages)
	return result
}

func (s storagePlugin) BackupDelete(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	backupDir := util.GetBackupDirFromConfig(config)
	backups, err := pluginUtil.ListBackups(backupDir)
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

				backupPath := backupDir + "/" + backupName
				err := pluginUtil.RecursiveDirDelete(backupPath)
				if err != nil {
					msg := util.SetMessage("ERROR", "Backup "+backupName+" delete failed! "+err.Error())
					messages = append(messages, msg)
					result = util.SetResult(1, messages)
					return result
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

func (s storagePlugin) BackupList(config util.Config) util.Backups {
	var backups util.Backups
	var result util.Result
	var messages []util.Message

	backupDir := util.GetBackupDirFromConfig(config)
	backupList, err := pluginUtil.ListBackups(backupDir)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)
		backups.Result = result

		return backups
	}

	result = util.SetResult(0, messages)
	backups.Result = result
	backups.Backups = backupList

	return backups
}

func (s storagePlugin) Mount(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "Mount not implemented or required")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) Unmount(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "Unmount not implemented or required")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "container-basic"
	plugin.Description = "Container Backup Plugin that uses rsync to backup a pod"
	plugin.Version = "1.0.0"
	plugin.Type = "storage"

	var capabilities []util.Capability
	var backupCap util.Capability
	backupCap.Name = "backup"

	var backupListCap util.Capability
	backupListCap.Name = "backupList"

	var backupDeleteCap util.Capability
	backupDeleteCap.Name = "backupDelete"

	var restoreCap util.Capability
	restoreCap.Name = "restore"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities, backupCap, backupListCap, backupDeleteCap, restoreCap, infoCap)

	plugin.Capabilities = capabilities

	return plugin
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

func main() {}
