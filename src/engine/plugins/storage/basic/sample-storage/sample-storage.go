package main

import (
	"os"
	"github.com/pborman/getopt/v2"
	"fossil/src/engine/util"
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
		fmt.Println("ERROR Incorrect parameter")
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
		fmt.Println("ERROR Incorrect parameter" + *optAction + "\n")
		getopt.Usage()
		os.Exit(1)
	}
}	

func backup (configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Backup ***")
}

func restore (configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Restore ***")
}

func backupList (configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Backup list ***")
}

func backupDelete (configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Backup delete ***")
}

func info () {
	var plugin util.Plugin = setPlugin()

	//output json
	b, err := json.Marshal(plugin)
    if err != nil {
        fmt.Println("ERROR " + err.Error())
	} else {
		fmt.Println(string(b))
	}
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "sample-storage"
	plugin.Description = "A sample storage plugin"
	plugin.Version = "1.0.0"
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
	config,err := util.ConfigMapToJson(configMap)
	if err != nil {
		fmt.Println("ERROR " + err.Error())
	}
	fmt.Println("DEBUG Config Parameters: " + config + "\n")
}

func getEnvParams() map[string]string {
	configMap := map[string]string{}

	configMap["ProfileName"] = os.Getenv("ProfileName")
	configMap["ConfigName"] = os.Getenv("ConfigName")
	configMap["BackupName"] = os.Getenv("BackupName")
	configMap["SelectedWorkflowId"] = os.Getenv("SelectedWorkflowId")
	configMap["AutoDiscovery"] = os.Getenv("AutoDiscovery")
	configMap["DataFilePaths"] = os.Getenv("DataFilePaths")
	configMap["LogFilePaths"] = os.Getenv("LogFilePaths")
	configMap["BackupPolicy"] = os.Getenv("BackupPolicy")
	configMap["SampleStorageVar1"] = os.Getenv("SampleStorageVar1")
	configMap["SampleStorageVar2"] = os.Getenv("SampleStorageVar2")

	return configMap
}