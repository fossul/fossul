package main

import (
	"os"
	"github.com/pborman/getopt/v2"
	"engine/util"
	"engine/client/k8s"
	"engine/client"
	"engine/plugins/pluginUtil"
	"encoding/json"
	"fmt"
)

func main() {
	optAction := getopt.StringLong("action",'a',"","backup|backupList|backupDelete|info")
	optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	if getopt.IsSet("action") != true {
		pluginUtil.LogErrorMessage("missing parameter --action")
		getopt.Usage()
		os.Exit(1)
	}

	//load env parameters
	configMap := getEnvParams()

	if *optAction == "backup" {
		backup(configMap)
	} else if *optAction == "backupList" {
		backupList(configMap)
	} else if *optAction == "backupDelete" {
		backupDelete(configMap)		
	} else if *optAction == "info" {
		info()			
	} else {
		pluginUtil.LogErrorMessage("incorrect parameter" + *optAction)
		getopt.Usage()
		os.Exit(1)
	}
}	

func backup (configMap map[string]string) {
	configMap = setWorkflowId(configMap)
	printEnv(configMap)

	pluginUtil.LogInfoMessage("Performing container backup")

	podName := k8s.GetPod(configMap["Namespace"],configMap["ServiceName"],configMap["AccessWithinCluster"])
	pluginUtil.LogErrorMessage("Performing backup for pod" + podName)

	backupName := util.GetBackupName(configMap["BackupName"],configMap["BackupPolicy"],configMap["WorkflowId"])
	backupPath := util.GetBackupPath(configMap)
	pluginUtil.LogInfoMessage("Backup name is " + backupName + ", Backup path is " + backupPath)

	pluginUtil.CreateDir(backupPath,0755)

	var args []string
	args = append(args,configMap["CopyCmdPath"])
	if configMap["ContainerPlatform"] == "openshift" {
		args = append(args,"rsync")
		args = append(args,"-n")
		args = append(args,configMap["Namespace"])
		args = append(args,podName + ":" + configMap["BackupSrcPath"])
		} else if configMap["ContainerPlatform"] == "kubernetes" {
			args = append(args,"cp")
			args = append(args,configMap["Namespace"] + "/" + podName + ":" + configMap["BackupSrcPath"])
	} else {

	}	
	args = append(args,backupPath)
	
	result := util.ExecuteCommand(args...)
	pluginUtil.LogResultMessages(result)
			
	if result.Code != 0 {
		os.Exit(1)
	}
}

func backupList (configMap map[string]string) {
	backupDir := util.GetBackupDir(configMap)
	backups := pluginUtil.ListBackups(backupDir)

	b, err := json.Marshal(backups)
    if err != nil {
        pluginUtil.LogErrorMessage(err.Error())
	} else {
		pluginUtil.PrintMessage(string(b))
	}
}

func backupDelete (configMap map[string]string) {
	printEnv(configMap)

	backupDir := util.GetBackupDir(configMap)
	backups := pluginUtil.ListBackups(backupDir)
	backupsByPolicy := util.GetBackupsByPolicy(configMap["BackupPolicy"],backups)
	backupCount := len(backupsByPolicy)
	backupRetentionCount := util.StringToInt(configMap["BackupRetention"])

	if backupCount > backupRetentionCount {
		count := 1
		for backup := range pluginUtil.ReverseBackupList(backupsByPolicy) {
			if count > backupRetentionCount {
				msg := fmt.Sprintf("Number of backups [%d] greater than backup retention [%d]",backupCount,backupRetentionCount)
				pluginUtil.LogInfoMessage(msg)
				backupCount = backupCount - 1

				backupName := backup.Name + "_" + backup.Policy + "_" + backup.WorkflowId + "_" + util.IntToString(backup.Epoch)
				pluginUtil.LogInfoMessage("Deleting backup " + backupName)
				backupPath := backupDir + "/" + backupName
				pluginUtil.RecursiveDirDelete(backupPath)
				pluginUtil.LogInfoMessage("Backup " + backupName + " deleted successfully")
				
				pluginUtil.LogInfoMessage("workflowId " + backup.WorkflowId)
				result := client.DeleteWorkflowResults(configMap["ProfileName"],configMap["ConfigName"],backup.WorkflowId)
				pluginUtil.LogResultMessages(result)

				if result.Code != 0 {
					os.Exit(1)
				}
			}
			count = count + 1
		}
	} else {
		msg := fmt.Sprintf("Backup deletion skipped, there are [%d] backups but backup retention is [%d]",backupCount, backupRetentionCount)
		pluginUtil.LogInfoMessage(msg)
	}
}

func info () {
	var plugin util.Plugin = setPlugin()

	b, err := json.Marshal(plugin)
    if err != nil {
        pluginUtil.LogErrorMessage(err.Error())
	} else {
		pluginUtil.PrintMessage(string(b))
	}
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "openshift-rsync"
	plugin.Description = "OpenShift Backup Plugin that uses rsync to backup a pod"
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

func printEnv(configMap map[string]string) {
	config := util.ConfigMapToJson(configMap)
	pluginUtil.LogDebugMessage("Config Parameters: " + config + "\n")
}

func getEnvParams() map[string]string {
	configMap := map[string]string{}

	configMap["ProfileName"] = os.Getenv("ProfileName")
	configMap["ConfigName"] = os.Getenv("ConfigName")
	configMap["BackupPolicy"] = os.Getenv("BackupPolicy")
	configMap["BackupRetention"] = os.Getenv("BackupRetention")
	configMap["BackupName"] = os.Getenv("BackupName")
	configMap["ContainerPlatform"] = os.Getenv("ContainerPlatform")
	configMap["AccessWithinCluster"] = os.Getenv("AccessWithinCluster")
	configMap["Namespace"] = os.Getenv("Namespace")
	configMap["ServiceName"] = os.Getenv("ServiceName")
	configMap["CopyCmdPath"] = os.Getenv("CopyCmdPath")
	configMap["BackupSrcPath"] = os.Getenv("BackupSrcPath")
	configMap["BackupDestPath"] = os.Getenv("BackupDestPath")

	return configMap
}

func setWorkflowId(configMap map[string]string) map[string]string {
	configMap["WorkflowId"] = os.Getenv("WorkflowId")

	return configMap
}