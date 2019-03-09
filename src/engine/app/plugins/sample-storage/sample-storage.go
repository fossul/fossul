package main

import (
	"os"
	"fmt"
	"github.com/pborman/getopt/v2"
	"engine/util"
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
		fmt.Println("ERROR: incorrect parameter")
		getopt.Usage()
		os.Exit(1)
	}

	if *optAction == "backup" {
		backup()
	} else if *optAction == "backupList" {
		backupList()
	} else if *optAction == "backupDelete" {
		backupDelete()		
	} else if *optAction == "info" {
		info()			
	} else {
		fmt.Println("ERROR: incorrect parameter", *optAction)
		getopt.Usage()
		os.Exit(1)
	}
}	

func backup () {
	printEnv()
	fmt.Println("Performing backup")
}

func backupList () {
	printEnv()
	fmt.Println("Performing backup list")
}

func backupDelete () {
	printEnv()
	fmt.Println("Performing backup delete")
}

func info () {
	var plugin util.Plugin = setPlugin()

	//output json
	b, err := json.Marshal(plugin)
    if err != nil {
        fmt.Println(err)
        return
	}
	
	fmt.Println(string(b))
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

func printEnv() {
	fmt.Println("Config Parameters")
	sampleStorageVar1 := os.Getenv("SampleStorageVar1")
	fmt.Println("SampleStorageVar1=" + sampleStorageVar1)

	sampleStorageVar2 := os.Getenv("SampleStorageVar2")
	fmt.Println("SampleStorageVar2=" + sampleStorageVar2)
}