package main

import (
	"os"
	"fmt"
	"github.com/pborman/getopt/v2"
	"engine/util"
	"engine/util/k8s"
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
		fmt.Printf("ERROR missing parameter --action")
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
		fmt.Printf("ERROR incorrect parameter" + *optAction)
		getopt.Usage()
		os.Exit(1)
	}
}	

func backup (configMap map[string]string) {
	printEnv(configMap)
	fmt.Println("INFO Performing container backup")

	podName := k8s.GetPod(configMap["Namespace"],configMap["ServiceName"],configMap["AccessWithinCluster"])
	fmt.Println("INFO Performing backup for pod" + podName)

	result := util.ExecuteCommand(configMap["RsyncCmdPath"],"rsync",podName + ":" + configMap["BackupSrcPath"],configMap["BackupDestPath"])

	for _, line := range result.Messages {
		//	//t := time.Unix(line.Timestamp,0)
		fmt.Println(line.Level, line.Message)
	}
			
	if result.Code != 0 {
		os.Exit(1)
	}
}

func backupList (configMap map[string]string) {
	printEnv(configMap)
	fmt.Printf("INFO Performing backup list")
}

func backupDelete (configMap map[string]string) {
	printEnv(configMap)
	fmt.Printf("INFO Performing backup delete")


}

func info () {
	var plugin util.Plugin = setPlugin()

	//output json
	b, err := json.Marshal(plugin)
    if err != nil {
        fmt.Println(err)
        return
	}
	
	fmt.Printf(string(b))
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
	fmt.Printf("DEBUG Config Parameters: " + config + "\n")
}

func getEnvParams() map[string]string {
	configMap := map[string]string{}

	configMap["ProfileName"] = os.Getenv("ProfileName")
	configMap["ConfigName"] = os.Getenv("ConfigName")
	configMap["AccessWithinCluster"] = os.Getenv("AccessWithinCluster")
	configMap["Namespace"] = os.Getenv("Namespace")
	configMap["ServiceName"] = os.Getenv("ServiceName")
	configMap["RsyncCmdPath"] = os.Getenv("RsyncCmdPath")
	configMap["BackupSrcPath"] = os.Getenv("BackupSrcPath")
	configMap["BackupDestPath"] = os.Getenv("BackupDestPath")

	return configMap
}