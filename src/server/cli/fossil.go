package main

import (
	"github.com/pborman/getopt/v2"
	"os"
	"fmt"
	"engine/util"
	
)

func main() {
    optAction := getopt.StringLong("action",'a',"","backup|list|status")
    optHelp := getopt.BoolLong("help", 0, "Help")
    getopt.Parse()

    if *optHelp {
        getopt.Usage()
        os.Exit(0)
    }
	
	if *optAction == "backup" {

	} else if *optAction == "list" {

	} else if *optAction == "status" {
		fmt.Println("Checking status of services")

		var workflowStatus util.Status
		workflowStatus = util.GetWorkflowServiceStatus()
		fmt.Println("Workflow Service:", util.Status(workflowStatus))

		var appStatus util.Status
		appStatus = util.GetAppServiceStatus()
		fmt.Println("App Service:", util.Status(appStatus))

		var storageStatus util.Status
		storageStatus = util.GetStorageServiceStatus()
		fmt.Println("Storage Service:", util.Status(storageStatus))

	}
}