package main

import (
	"engine/util"
	"engine/client/k8s"
	"engine/plugins/pluginUtil"
	"engine/client"
	"log"
	"fmt"
)

type storagePlugin string

var config util.Config
var StoragePlugin storagePlugin

func (s storagePlugin) SetEnv(c util.Config) util.Result {
	config = c
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	result = util.SetResult(resultCode, messages)

	return result
}

func (s storagePlugin) Backup() util.Result {	
	var result util.Result
	var messages []util.Message

	msg := util.SetMessage("INFO", "Performing container backup")
	messages = append(messages,msg)

	podName,err := k8s.GetPod(config.StoragePluginParameters["Namespace"],config.StoragePluginParameters["ServiceName"],config.StoragePluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages,msg)

		result = util.SetResult(1, messages)
		return result
	}

	msg = util.SetMessage("INFO", "Performing backup for pod " + podName)
	messages = append(messages,msg)

	backupName := util.GetBackupName(config.StoragePluginParameters["BackupName"],config.SelectedBackupPolicy,config.WorkflowId)
	backupPath := util.GetBackupPathFromConfig(config)
	msg = util.SetMessage("INFO", "Backup name is " + backupName + ", Backup path is " + backupPath)
	messages = append(messages,msg)

	err = pluginUtil.CreateDir(backupPath,0755)
	if err != nil {
		msg = util.SetMessage("ERROR", err.Error())
		messages = append(messages,msg)
		result = util.SetResult(1, messages)
		return result
	}

	var args []string
	args = append(args,config.StoragePluginParameters["CopyCmdPath"])
	if config.StoragePluginParameters["ContainerPlatform"] == "openshift" {
		args = append(args,"rsync")
		args = append(args,"-n")
		args = append(args,config.StoragePluginParameters["Namespace"])
		args = append(args,podName + ":" + config.StoragePluginParameters["BackupSrcPath"])
	} else if config.StoragePluginParameters["ContainerPlatform"] == "kubernetes" {
		args = append(args,"cp")
		args = append(args,config.StoragePluginParameters["Namespace"] + "/" + podName + ":" + config.StoragePluginParameters["BackupSrcPath"])
	} else {

	}	
	args = append(args,backupPath)
	
	cmdResult := util.ExecuteCommand(args...)
	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		for _,msg := range cmdResult.Messages {
			messages = util.PrependMessage(msg,messages)
		}
	}

	result = util.SetResult(0, messages)
	return result
}

func (s storagePlugin) BackupDelete() util.Result {	
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	backupDir := util.GetBackupDirFromConfig(config)
	backups,err := pluginUtil.ListBackups(backupDir)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages,msg)
		result = util.SetResult(1, messages)
		return result		
	}

	backupsByPolicy := util.GetBackupsByPolicy(config.SelectedBackupPolicy,backups)
	backupCount := len(backupsByPolicy)

	if backupCount > config.SelectedBackupRetention {
		count := 1
		for backup := range pluginUtil.ReverseBackupList(backupsByPolicy) {
			if count > config.SelectedBackupRetention {
				msg := util.SetMessage("INFO",fmt.Sprintf("Number of backups [%d] greater than backup retention [%d]",backupCount,config.SelectedBackupRetention))
				messages = append(messages,msg)
				backupCount = backupCount - 1

				backupName := backup.Name + "_" + backup.Policy + "_" + backup.WorkflowId + "_" + util.IntToString(backup.Epoch)
				msg = util.SetMessage("INFO", "Deleting backup " + backupName)
				messages = append(messages,msg)
	
				backupPath := backupDir + "/" + backupName
				pluginUtil.RecursiveDirDelete(backupPath)
				msg = util.SetMessage("INFO", "Backup " + backupName + " deleted successfully")
				messages = append(messages,msg)

				deleteWorkflowResult := client.DeleteWorkflowResults(config.ProfileName,config.ConfigName,backup.WorkflowId)
				if deleteWorkflowResult.Code != 0 {
					for _,msg := range deleteWorkflowResult.Messages {
						messages = append(messages,msg)
					}
					result = util.SetResult(1, messages)
					return result
				} else {
					for _,msg := range deleteWorkflowResult.Messages {
						messages = append(messages,msg)
					}
				}
			}
			count = count + 1
		}
	} else {
		msg := util.SetMessage("INFO",fmt.Sprintf("Backup deletion skipped, there are [%d] backups but backup retention is [%d]",backupCount, config.SelectedBackupRetention))
		messages = append(messages,msg)
	}	

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) BackupList() []util.Backup {	
	var backups []util.Backup
	var err error

	backupDir := util.GetBackupDirFromConfig(config)
	backups,err = pluginUtil.ListBackups(backupDir)
	if err != nil {
		log.Println(err.Error())
	}

	return backups
}

func (s storagePlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "container-basic"
	plugin.Description = "Container Backup Plugin that uses rsync to backup a pod"
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

	capabilities = append(capabilities,backupCap,backupListCap,backupDeleteCap,infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func main() {}