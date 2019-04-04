package main

import (
	"os"
	"github.com/pborman/getopt/v2"
	"engine/util"
	"encoding/json"
	"fmt"
)

func main() {
	optAction := getopt.StringLong("action",'a',"","archive|archiveList|archiveDelete")
	optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	if getopt.IsSet("action") != true {
		fmt.Println("ERROR Incorrect parameter\n")
		getopt.Usage()
		os.Exit(1)
	}

		//load env parameters
		configMap := getEnvParams()

	if *optAction == "archive" {
		archive(configMap)
	} else if *optAction == "archiveList" {
		archiveList(configMap)
	} else if *optAction == "archiveDelete" {
		archiveDelete(configMap)		
	} else if *optAction == "info" {
		info()			
	} else {
		fmt.Println("ERROR Incorrect parameter" + *optAction + "\n")
		getopt.Usage()
		os.Exit(1)
	}
}	

func archive (configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Archive ***")
}

func archiveList (configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Archive list ***")
}

func archiveDelete (configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Archive delete ***")
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
	plugin.Name = "sample-archive"
	plugin.Description = "A sample archive plugin"
	plugin.Type = "archive"

	var capabilities []util.Capability
	var archiveCap util.Capability
	archiveCap.Name = "archive"

	var archiveListCap util.Capability
	archiveListCap.Name = "archiveList"

	var archiveDeleteCap util.Capability
	archiveDeleteCap.Name = "archiveDelete"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities,archiveCap,archiveListCap,archiveDeleteCap,infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func printEnv(configMap map[string]string) {
	config := util.ConfigMapToJson(configMap)
	fmt.Println("DEBUG Config Parameters: " + config + "\n")
}

func getEnvParams() map[string]string {
	configMap := map[string]string{}

	configMap["ProfileName"] = os.Getenv("ProfileName")
	configMap["ConfigName"] = os.Getenv("ConfigName")
	configMap["BackupName"] = os.Getenv("BackupName")
	configMap["SampleArchiveVar1"] = os.Getenv("SampleArchiveVar1")
	configMap["SampleArchiveVar2"] = os.Getenv("SampleArchiveVar2")

	return configMap
}