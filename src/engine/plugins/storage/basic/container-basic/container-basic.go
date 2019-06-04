package main

import (
	"encoding/json"
	"fmt"
	"fossul/src/engine/client/k8s"
	"fossul/src/engine/plugins/pluginUtil"
	"fossul/src/engine/util"
	"github.com/pborman/getopt/v2"
	"os"
	"strings"
)

func main() {
	optAction := getopt.StringLong("action", 'a', "", "backup|restore|backupList|backupDelete|info")
	optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	if getopt.IsSet("action") != true {
		fmt.Println("ERROR missing parameter --action")
		getopt.Usage()
		os.Exit(1)
	}

	//load env parameters
	configMap := getEnvParams()

	if *optAction == "backup" {
		backup(configMap)
	} else if *optAction == "restore" {
		restore(configMap)
	} else if *optAction == "backupList" {
		backupList(configMap)
	} else if *optAction == "backupDelete" {
		backupDelete(configMap)
	} else if *optAction == "info" {
		info()
	} else {
		fmt.Println("ERROR incorrect parameter" + *optAction)
		getopt.Usage()
		os.Exit(1)
	}
}

func backup(configMap map[string]string) {
	configMap = setWorkflowId(configMap)
	printEnv(configMap)

	backupSrcFilePaths := getBackupSrcPaths(configMap)

	fmt.Println("INFO Performing container backup")

	podName, err := k8s.GetPod(configMap["Namespace"], configMap["ServiceName"], configMap["AccessWithinCluster"])
	checkError(err)

	fmt.Println("INFO Performing backup for pod " + podName)

	backupName := util.GetBackupName(configMap["BackupName"], configMap["BackupPolicy"], configMap["WorkflowId"])
	backupPath := util.GetBackupPathFromMap(configMap)
	fmt.Println("INFO Backup name is " + backupName + ", Backup path is " + backupPath)

	err = pluginUtil.CreateDir(backupPath, 0755)
	checkError(err)

	for _, backupSrcFilePath := range backupSrcFilePaths {
		var args []string
		args = append(args, configMap["CopyCmdPath"])
		if configMap["ContainerPlatform"] == "openshift" {
			args = append(args, "rsync")
			args = append(args, "-n")
			args = append(args, configMap["Namespace"])
			args = append(args, podName+":"+backupSrcFilePath)
		} else if configMap["ContainerPlatform"] == "kubernetes" {
			args = append(args, "cp")
			args = append(args, configMap["Namespace"]+"/"+podName+":"+backupSrcFilePath)
		} else {
			fmt.Println("ERROR incorrect parameter set for ContainerPlatform [" + configMap["ContainerPlatform"] + "]")
			os.Exit(1)
		}
		args = append(args, backupPath)

		result := util.ExecuteCommand(args...)
		for _, line := range result.Messages {
			fmt.Println(line.Level, line.Message)
		}

		if result.Code != 0 {
			os.Exit(1)
		}
	}
}

func restore(configMap map[string]string) {
	configMap = setWorkflowId(configMap)
	printEnv(configMap)

	fmt.Println("INFO Performing container restore")

	podName, err := k8s.GetPod(configMap["Namespace"], configMap["ServiceName"], configMap["AccessWithinCluster"])
	checkError(err)

	fmt.Println("INFO Performing restore for pod " + podName)

	restorePath, err := util.GetRestoreSrcPathFromMap(configMap)
	checkError(err)

	if restorePath == "" {
		fmt.Println("ERROR Restore data no longer available for workflow id [" + configMap["SelectedWorkflowId"] + "], check retention policy")
		os.Exit(1)
	}

	fmt.Println("INFO Restore source path is [" + restorePath + "]")

	restoreDestPath := "/tmp/" + configMap["SelectedWorkflowId"]

	var args []string
	args = append(args, configMap["CopyCmdPath"])
	if configMap["ContainerPlatform"] == "openshift" {
		args = append(args, "rsync")
		args = append(args, "-n")
		args = append(args, configMap["Namespace"])
		args = append(args, restorePath)
		args = append(args, podName+":"+restoreDestPath)
	} else if configMap["ContainerPlatform"] == "kubernetes" {
		args = append(args, "cp")
		args = append(args, restorePath)
		args = append(args, configMap["Namespace"]+"/"+podName+":"+restoreDestPath)
	} else {
		fmt.Println("ERROR incorrect parameter set for ContainerPlatform [" + configMap["ContainerPlatform"] + "]")
		os.Exit(1)
	}

	result := util.ExecuteCommand(args...)
	for _, line := range result.Messages {
		fmt.Println(line.Level, line.Message)
	}

	if result.Code != 0 {
		os.Exit(1)
	}
}

func backupList(configMap map[string]string) {
	backupDir := util.GetBackupDirFromMap(configMap)
	backups, err := pluginUtil.ListBackups(backupDir)
	checkError(err)

	b, err := json.Marshal(backups)
	if err != nil {
		fmt.Println("ERROR " + err.Error())
	} else {
		fmt.Println(string(b))
	}
}

func backupDelete(configMap map[string]string) {
	printEnv(configMap)

	backupDir := util.GetBackupDirFromMap(configMap)
	backups, err := pluginUtil.ListBackups(backupDir)
	checkError(err)

	backupsByPolicy := util.GetBackupsByPolicy(configMap["BackupPolicy"], backups)
	backupCount := len(backupsByPolicy)
	backupRetentionCount := util.StringToInt(configMap["BackupRetention"])

	if backupCount > backupRetentionCount {
		count := 1
		for backup := range pluginUtil.ReverseBackupList(backupsByPolicy) {
			if count > backupRetentionCount {
				msg := fmt.Sprintf("Number of backups [%d] greater than backup retention [%d]", backupCount, backupRetentionCount)
				fmt.Println(msg)
				backupCount = backupCount - 1

				backupName := backup.Name + "_" + backup.Policy + "_" + backup.WorkflowId + "_" + util.IntToString(backup.Epoch)
				fmt.Println("INFO Deleting backup " + backupName)
				backupPath := backupDir + "/" + backupName
				err := pluginUtil.RecursiveDirDelete(backupPath)
				if err != nil {
					fmt.Println("ERROR Backup " + backupName + " delete failed! " + err.Error())
					os.Exit(1)
				}
				fmt.Println("INFO Backup " + backupName + " deleted successfully")
			}
			count = count + 1
		}
	} else {
		msg := fmt.Sprintf("INFO Backup deletion skipped, there are [%d] backups but backup retention is [%d]", backupCount, backupRetentionCount)
		fmt.Println(msg)
	}
}

func info() {
	var plugin util.Plugin = setPlugin()

	b, err := json.Marshal(plugin)
	if err != nil {
		fmt.Println("ERROR " + err.Error())
	} else {
		fmt.Println(string(b))
	}
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

func printEnv(configMap map[string]string) {
	config, err := util.ConfigMapToJson(configMap)
	if err != nil {
		fmt.Println("ERROR " + err.Error())
	}
	fmt.Println("DEBUG Config Parameters: " + config + "\n")
}

func getEnvParams() map[string]string {
	configMap := map[string]string{}

	configMap["ProfileName"] = os.Getenv("ProfileName")
	configMap["ConfigName"] = os.Getenv("ConfigName")
	configMap["SelectedWorkflowId"] = os.Getenv("SelectedWorkflowId")
	configMap["AutoDiscovery"] = os.Getenv("AutoDiscovery")
	configMap["DataFilePaths"] = os.Getenv("DataFilePaths")
	configMap["LogFilePaths"] = os.Getenv("LogFilePaths")
	configMap["BackupPolicy"] = os.Getenv("BackupPolicy")
	configMap["BackupRetention"] = os.Getenv("BackupRetention")
	configMap["BackupName"] = os.Getenv("BackupName")
	configMap["ContainerPlatform"] = os.Getenv("ContainerPlatform")
	configMap["AccessWithinCluster"] = os.Getenv("AccessWithinCluster")
	configMap["Namespace"] = os.Getenv("Namespace")
	configMap["ServiceName"] = os.Getenv("ServiceName")
	configMap["CopyCmdPath"] = os.Getenv("CopyCmdPath")
	configMap["BackupSrcPaths"] = os.Getenv("BackupSrcPaths")
	configMap["BackupDestPath"] = os.Getenv("BackupDestPath")

	return configMap
}

func setWorkflowId(configMap map[string]string) map[string]string {
	configMap["WorkflowId"] = os.Getenv("WorkflowId")

	return configMap
}

func checkError(err error) {
	if err != nil {
		fmt.Println("ERROR " + err.Error())
		os.Exit(1)
	}
}

func getBackupSrcPaths(configMap map[string]string) []string {
	var backupSrcFilePaths []string
	if configMap["AutoDiscovery"] == "true" {
		dataPaths := strings.Split(configMap["DataFilePaths"], ",")
		logPaths := strings.Split(configMap["LogFilePaths"], ",")

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
		backupSrcFilePaths = strings.Split(configMap["BackupSrcPaths"], ",")
	}

	return backupSrcFilePaths
}
