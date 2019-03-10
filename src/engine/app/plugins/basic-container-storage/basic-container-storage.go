package main

import (
	"os"
	"fmt"
	"github.com/pborman/getopt/v2"
	"engine/util"
	"encoding/json"
	"engine/util/k8s"
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
		fmt.Printf("ERROR: incorrect parameter\n")
		getopt.Usage()
		os.Exit(1)
	}

	//load env parameters
	configMap := getEnvParams()
	printEnv(configMap)

	if *optAction == "backup" {
		backup(configMap)
	} else if *optAction == "backupList" {
		backupList(configMap)
	} else if *optAction == "backupDelete" {
		backupDelete(configMap)		
	} else if *optAction == "info" {
		info()			
	} else {
		fmt.Printf("ERROR: incorrect parameter" + *optAction + "\n")
		getopt.Usage()
		os.Exit(1)
	}
}	

func backup (configMap map[string]string) {
	fmt.Printf("Performing container backup")

	podName := k8s.GetPod(configMap["Namespace"],configMap["ServiceName"],configMap["AccessWithinCluster"])
	fmt.Printf("INFO: Performing backup for pod" + podName)
}

func backupList (configMap map[string]string) {
	fmt.Printf("Performing backup list")
}

func backupDelete (configMap map[string]string) {
	fmt.Printf("Performing backup delete")
}

func info () {
	var plugin util.Plugin = setPlugin()

	//output json
	b, err := json.Marshal(plugin)
    if err != nil {
        fmt.Println(err)
        return
	}
	
	fmt.Printf(string(b) + "\n")
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
	fmt.Printf("Config Parameters\n")
	fmt.Printf("AccessWithinCluster=" + configMap["AccessWithinCluster"] + "\n")
	fmt.Printf("Namespace=" + configMap["Namespace"] + "\n")
	fmt.Printf("SeriveName=" + configMap["ServiceName"] + "\n")
}

func getEnvParams() map[string]string {
	configMap := map[string]string{}

	configMap["AccessWithinCluster"] = os.Getenv("AccessWithinCluster")
	configMap["Namespace"] = os.Getenv("Namespace")
	configMap["ServiceName"] = os.Getenv("ServiceName")

	return configMap
}