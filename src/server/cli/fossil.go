package main

import (
	"github.com/pborman/getopt/v2"
	"os"
	"engine/util"
	"engine/client"
	"fmt"
//	"time"
)

func main() {
	optProfile := getopt.StringLong("profile",'p',"","Profile name")
	optConfig := getopt.StringLong("config",'c',"","Config name")
	optConfigPath := getopt.StringLong("configPath",'o',"","Path to configs directory")
	optPolicy := getopt.StringLong("policy",'i',"","Backup policy as defined in config")
	optAction := getopt.StringLong("action",'a',"","backup|backupList|appPluginList|storagePluginList|pluginInfo|status")
	optPluginName := getopt.StringLong("plugin",'l',"","Name of plugin")
    optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

    if *optHelp {
        getopt.Usage()
        os.Exit(0)
	}
	if getopt.IsSet("profile") != true {
		fmt.Println("ERROR: missing parameter --profile")
		getopt.Usage()
		os.Exit(1)	
	} else if getopt.IsSet("config") != true {
		fmt.Println("ERROR: missing parameter --config")
		getopt.Usage()
		os.Exit(1)
	} else if getopt.IsSet("action") != true {
		fmt.Println("ERROR: missing parameter --action")
		getopt.Usage()
		os.Exit(1)
	}

	var configPath string
	if getopt.IsSet("configPath") == true {
		configPath = *optConfigPath + "/" + *optProfile + "/" + *optConfig + ".conf"
	} else {
		configPath = "configs/" + *optProfile + "/" + *optConfig
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println(err,"\n" + "ERROR: Profile of Config don't exist")
		os.Exit(1)
	}

	//read config file into struct
	var config util.Config = util.ReadConfig(configPath)

	// Check retention policy
	if *optAction == "backup" || *optAction == "backupList" {
		if getopt.IsSet("policy") != true {
			fmt.Println("ERROR: missing parameter --policy")
			getopt.Usage()
			os.Exit(1)	
		}
		if util.ExistsBackupRetention(*optPolicy,config.BackupRetentions) != true {
			fmt.Println("ERROR: policy [" + *optPolicy + "] does npot match policy defined in config")
			os.Exit(1)
		}	
	}

	//load dynamic plugin parameters into config struct
	appConfigPath := *optConfigPath + "/" + *optProfile + "/" + config.AppPlugin + ".conf"
	storageConfigPath := *optConfigPath + "/" + *optProfile + "/" + config.StoragePlugin + ".conf"
	config = util.SetPluginParameters(appConfigPath, storageConfigPath, config)

	if *optAction == "backup" {
		logger := util.GetLoggerInstance()
		//util.LogCommentMessage(logger, "Welcome To Fossil Backup Framework, Performing Backup")

		var result []util.Result
		result = client.StartBackupWorkflow(string(*optProfile),string(*optConfig),string(*optPolicy),config)

		util.LogResults(logger, result)

		/*
		for _, item := range result {
			for _, line := range item.Messages {
				//t := time.Unix(line.Timestamp,0)
				//fmt.Println(t.Format("2006-01-02 15:04:05"), line.Level, line.Message)
				fmt.Println(time.Unix(line.Timestamp,0), line.Level, line.Message)
			}
			//fmt.Println(index, item.Time, item.Code, item.Messages)
		}
		*/	


	} else if *optAction == "backupList" {
		fmt.Println("### List of Backups ###")
		result, backups := client.BackupList(string(*optProfile),string(*optConfig),string(*optPolicy),config)
		backupsByPolicy := util.GetBackupsByPolicy(string(*optPolicy),backups)

		checkResult(result)

		for _, backup := range backupsByPolicy {
			fmt.Println(backup.Name, backup.Policy, backup.Timestamp)
		}
	} else if *optAction == "appPluginList" {
		fmt.Println("### List of Application Plugins ###")

		var plugins []string
		var appPlugins []string = client.AppPluginList(config)
		plugins = util.JoinArray(appPlugins,plugins)

		for _, plugin := range plugins {
			fmt.Println(plugin)
		}
	} else if *optAction == "storagePluginList" {
		fmt.Println("### List of Storage Plugins ###")
	
		var plugins []string
		var storagePlugins []string = client.StoragePluginList(config)
		plugins = util.JoinArray(storagePlugins,plugins)
	
		for _, plugin := range plugins {
			fmt.Println(plugin)
		}		
	} else if *optAction == "pluginInfo" {
		if getopt.IsSet("plugin") != true {
			fmt.Println("ERROR: Missing parameter --plugin")
			getopt.Usage()
			os.Exit(1)
		}
		var pluginName string = *optPluginName

		var result util.ResultSimple
		var plugin util.Plugin

		var appPlugins []string = client.AppPluginList(config)
		isAppPlugin := util.ExistsInArray(appPlugins,pluginName)

		if isAppPlugin == false {
			var storagePlugins []string = client.StoragePluginList(config)	
			isStoragePlugin := util.ExistsInArray(storagePlugins,pluginName)
			if isStoragePlugin == true {
				result, plugin = client.StoragePluginInfo(config,pluginName)
			} else {
				error := fmt.Sprintf("ERROR: Plugin %s not found!", plugin)
				fmt.Println(error)
			}
		} else if isAppPlugin == true {
			result, plugin = client.AppPluginInfo(config,pluginName)
		} else {
			error := fmt.Sprintf("ERROR: Plugin %s not found!", plugin)
			fmt.Println(error)
		}				

		checkResult(result)

		fmt.Println("### Plugin Information ###")
		fmt.Println("Name:", plugin.Name)
		fmt.Println("Description:", plugin.Description)
		fmt.Println("Type:", plugin.Type)
		fmt.Println("Capabilities:", plugin.Capabilities)
	} else if *optAction == "status" {
		fmt.Println("### Checking status of services ###")

		var workflowStatus util.Status
		workflowStatus = client.GetWorkflowServiceStatus()
		fmt.Println("Workflow Service:", workflowStatus)

		var appStatus util.Status
		appStatus = client.GetAppServiceStatus()
		fmt.Println("App Service:", appStatus)

		var storageStatus util.Status
		storageStatus = client.GetStorageServiceStatus()
		fmt.Println("Storage Service:", storageStatus)
	} else {
		fmt.Println("ERROR: incorrect parameter", *optAction)
		getopt.Usage()
		os.Exit(1)
	}
}

func checkResult(result util.ResultSimple) {
	if result.Code != 0 {
		for _, line := range result.Messages {
			fmt.Println(line)
		}
	}
}