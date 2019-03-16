package main

import (
	"os"
	"github.com/pborman/getopt/v2"
	"engine/util"
	"engine/client/k8s"
	"engine/plugins/pluginUtil"
	"encoding/json"
	"strconv"
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
	printEnv(configMap)
	pluginUtil.LogInfoMessage("Performing container backup")

	podName := k8s.GetPod(configMap["Namespace"],configMap["ServiceName"],configMap["AccessWithinCluster"])
	pluginUtil.LogErrorMessage("Performing backup for pod" + podName)

	backupName := util.GetBackupName(configMap["BackupName"],configMap["BackupPolicy"])
	backupPath := util.GetBackupPath(configMap)
	pluginUtil.LogInfoMessage("Backup name is " + backupName + ", Backup path is " + backupPath)

	pluginUtil.CreateDir(backupPath,0755)

	result := util.ExecuteCommand(configMap["RsyncCmdPath"],"rsync",podName + ":" + configMap["BackupSrcPath"],backupPath)
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
	backupRetentionCount, err := strconv.Atoi(configMap["BackupRetention"])
	if err != nil {
		pluginUtil.LogErrorMessage(err.Error())
	}


	if backupCount > backupRetentionCount {
		msg := fmt.Sprintf("Number of backups [%d] greater than backup retention [%d]",backupCount,backupRetentionCount)
		pluginUtil.LogInfoMessage(msg)
		count := 1
		for backup := range pluginUtil.ReverseBackupList(backupsByPolicy) {
			if count > backupRetentionCount {
				pluginUtil.LogInfoMessage("Deleting backup " + backup.Name + "_" + backup.Epoch)
				backupPath := backupDir + "/" + backup.Name + "_" + backup.Policy + "_" + backup.Epoch
				pluginUtil.RecursiveDirDelete(backupPath)
				pluginUtil.LogInfoMessage("Backup " + backup.Name + "_" + backup.Policy + "_" + backup.Epoch + " deleted successfully")
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
	configMap["AccessWithinCluster"] = os.Getenv("AccessWithinCluster")
	configMap["Namespace"] = os.Getenv("Namespace")
	configMap["ServiceName"] = os.Getenv("ServiceName")
	configMap["RsyncCmdPath"] = os.Getenv("RsyncCmdPath")
	configMap["BackupSrcPath"] = os.Getenv("BackupSrcPath")
	configMap["BackupDestPath"] = os.Getenv("BackupDestPath")

	return configMap
}