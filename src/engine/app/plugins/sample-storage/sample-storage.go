package main

import (
	"os"
	"github.com/pborman/getopt/v2"
	"engine/util"
	"engine/util/pluginUtil"
	"encoding/json"
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
		pluginUtil.LogErrorMessage("Incorrect parameter\n")
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
		pluginUtil.LogErrorMessage("Incorrect parameter" + *optAction + "\n")
		getopt.Usage()
		os.Exit(1)
	}
}	

func backup (configMap map[string]string) {
	printEnv(configMap)
	pluginUtil.LogInfoMessage("Performing backup")
}

func backupList (configMap map[string]string) {
	printEnv(configMap)
	pluginUtil.LogErrorMessage("Performing backup list")
}

func backupDelete (configMap map[string]string) {
	printEnv(configMap)
	pluginUtil.LogErrorMessage("Performing backup delete")
}

func info () {
	var plugin util.Plugin = setPlugin()

	//output json
	b, err := json.Marshal(plugin)
    if err != nil {
        pluginUtil.LogErrorMessage(err.Error())
	} else {
		pluginUtil.PrintMessage(string(b))
	}
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "sample"
	plugin.Description = "A sample plugin"
	plugin.Type = "app"

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
	configMap["BackupName"] = os.Getenv("BackupName")
	configMap["SampleStorageVar1"] = os.Getenv("SampleStorageVar1")
	configMap["SampleStorageVar2"] = os.Getenv("SampleStorageVar2")

	return configMap
}