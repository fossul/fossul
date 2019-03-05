package main

import (
	"github.com/pborman/getopt/v2"
	"os"
	"fmt"
	"engine/util"
	
)

func main() {
	optConfig := getopt.StringLong("config",'c',"","Path to config file")
    optAction := getopt.StringLong("action",'a',"","backup|list|status")
    optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

    if *optHelp {
        getopt.Usage()
        os.Exit(0)
	}
	
	configPath := *optConfig
	var config util.Config = util.ReadConfig(configPath)
	fmt.Println(config.BackupRetentions[1].Policy)

	if *optAction == "backup" {

		var result []util.Result
		result = util.StartBackupWorkflow(config)

		for index, item := range result {
			fmt.Println(index, item.Time, item.Code, item.Stdout, item.Stderr)
		}	


	} else if *optAction == "list" {

	} else if *optAction == "status" {
		fmt.Println("Checking status of services")

		var workflowStatus util.Status
		workflowStatus = util.GetWorkflowServiceStatus()
		fmt.Println("Workflow Service:", workflowStatus)

		var appStatus util.Status
		appStatus = util.GetAppServiceStatus()
		fmt.Println("App Service:", appStatus)

		var storageStatus util.Status
		storageStatus = util.GetStorageServiceStatus()
		fmt.Println("Storage Service:", storageStatus)

	}
}