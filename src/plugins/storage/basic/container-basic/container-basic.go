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
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fossul/fossul/src/client/k8s"
	"github.com/fossul/fossul/src/engine/util"
	"github.com/fossul/fossul/src/plugins/pluginUtil"
	"github.com/pborman/getopt/v2"
)

func main() {
	optBackup := getopt.BoolLong("backup", 0, "Backup")
	optRestore := getopt.BoolLong("restore", 0, "Restore")
	optBackupList := getopt.BoolLong("backupList", 0, "Backup List")
	optBackupDelete := getopt.BoolLong("backupDeleteWorkflow", 0, "Backup Delete")
	optMount := getopt.BoolLong("mount", 0, "Mount a backup or snapshot volume")
	optUnmount := getopt.BoolLong("unmount", 0, "Mount a backup or snapshot volume")
	optInfo := getopt.BoolLong("info", 0, "Storage Plugin Information")
	optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	//load env parameters
	configMap := getEnvParams()

	if *optBackup {
		backup(configMap)
	} else if *optRestore {
		restore(configMap)
	} else if *optBackupList {
		backupList(configMap)
	} else if *optBackupDelete {
		backupDeleteWorkflow(configMap)
	} else if *optMount {
		mount(configMap)
	} else if *optUnmount {
		unmount(configMap)
	} else if *optInfo {
		info()
	} else {
		getopt.Usage()
		os.Exit(0)
	}
}

func backup(configMap map[string]string) {
	configMap = setWorkflowId(configMap)
	printEnv(configMap)

	backupSrcFilePaths := getBackupSrcPaths(configMap)

	fmt.Println("INFO Performing container backup")

	podName, err := k8s.GetPodName(configMap["Namespace"], configMap["ServiceName"], configMap["AccessWithinCluster"])
	checkError(err)

	fmt.Println("INFO Performing backup for pod " + podName)

	backupName := util.GetBackupName(configMap["BackupName"], configMap["BackupPolicy"], configMap["WorkflowId"], configMap["WorkflowTimestamp"])
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

	podName, err := k8s.GetPodName(configMap["Namespace"], configMap["ServiceName"], configMap["AccessWithinCluster"])
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

func backupDeleteWorkflow(configMap map[string]string) {
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

				backupName := backup.Name + "-" + backup.Policy + "-" + backup.WorkflowId + "-" + util.IntToString(backup.Epoch)
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

func mount(configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Mount ***")
}

func unmount(configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Unmount ***")
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
	backupDeleteCap.Name = "backupDeleteWorkflow"

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
	configMap["WorkflowId"] = os.Getenv("WorkflowId")
	configMap["WorkflowTimestamp"] = os.Getenv("WorkflowTimestamp")
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
