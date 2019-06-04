package main

import (
	"encoding/json"
	"fmt"
	"fossul/src/engine/util"
	"github.com/pborman/getopt/v2"
	"os"
	"strings"
)

func main() {
	optAction := getopt.StringLong("action", 'a', "", "discover|quiesce|unquiesce|info")
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

	if *optAction == "quiesce" {
		printEnv(configMap)
		quiesce(configMap)
	} else if *optAction == "unquiesce" {
		printEnv(configMap)
		unquiesce(configMap)
	} else if *optAction == "preRestore" {
		printEnv(configMap)
		preRestore(configMap)
	} else if *optAction == "postRestore" {
		printEnv(configMap)
		postRestore(configMap)
	} else if *optAction == "info" {
		info()
	} else if *optAction == "discover" {
		discover()
	} else {
		fmt.Println("ERROR Incorrect parameter" + *optAction + "\n")
		getopt.Usage()
		os.Exit(1)
	}
}

func discover() {
	var discoverResult util.DiscoverResult = setDiscoverResult()

	//output json
	b, err := json.Marshal(discoverResult)
	if err != nil {
		fmt.Println("ERROR " + err.Error())
	} else {
		fmt.Println(string(b))
	}
}

func quiesce(configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Application quiesce ***")
}

func unquiesce(configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Application unquiesce ***")
}

func preRestore(configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Application Pre Restore ***")
}

func postRestore(configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO *** Application Post Restore ***")
}

func info() {
	var plugin util.Plugin = setPlugin()

	//output json
	b, err := json.Marshal(plugin)
	if err != nil {
		fmt.Println("ERROR " + err.Error())
	} else {
		fmt.Println(string(b))
	}
}

func setDiscoverResult() (discoverResult util.DiscoverResult) {
	var data []string
	data = append(data, "/path/to/data/file1")
	data = append(data, "/path/to/data/file2")

	var logs []string
	logs = append(logs, "/path/to/logs/file1")
	logs = append(logs, "/path/to/logs/file2")

	var discoverInst1 util.Discover
	discoverInst1.Instance = "inst1"
	discoverInst1.DataFilePaths = data
	discoverInst1.LogFilePaths = logs

	var discoverInst2 util.Discover
	discoverInst2.Instance = "inst2"
	discoverInst2.DataFilePaths = data
	discoverInst2.LogFilePaths = logs

	var discoverList []util.Discover
	discoverList = append(discoverList, discoverInst1)
	discoverList = append(discoverList, discoverInst2)

	var messages []util.Message
	msg := util.SetMessage("INFO", "*** Application Discovery ***")
	messages = append(messages, msg)

	for _, discover := range discoverList {
		dataFiles := strings.Join(discover.DataFilePaths, " ")
		logFiles := strings.Join(discover.LogFilePaths, " ")
		msg := util.SetMessage("INFO", "Instance ["+discover.Instance+"] data files: ["+dataFiles+"] log files: ["+logFiles+"]")
		messages = append(messages, msg)
	}

	result := util.SetResult(0, messages)
	discoverResult.Result = result
	discoverResult.DiscoverList = discoverList

	return discoverResult
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "sample"
	plugin.Description = "A sample plugin"
	plugin.Version = "1.0.0"
	plugin.Type = "app"

	var capabilities []util.Capability
	var discoverCap util.Capability
	discoverCap.Name = "discover"

	var quiesceCap util.Capability
	quiesceCap.Name = "quiesce"

	var unquiesceCap util.Capability
	unquiesceCap.Name = "unquiesce"

	var preRestoreCap util.Capability
	preRestoreCap.Name = "preRestore"

	var postRestoreCap util.Capability
	postRestoreCap.Name = "postRestore"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities, discoverCap, quiesceCap, unquiesceCap, preRestoreCap, postRestoreCap, infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func printEnv(configMap map[string]string) {
	config, err := util.ConfigMapToJson(configMap)
	if err != nil {
		fmt.Println("ERROR " + err.Error())
	}
	fmt.Println("DEBUG Config Parameters: " + config)
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
	configMap["SampleAppVar1"] = os.Getenv("SampleAppVar1")
	configMap["SampleAppVar2"] = os.Getenv("SampleAppVar2")

	return configMap
}
