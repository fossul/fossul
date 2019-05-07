package main

import (
	"github.com/pborman/getopt/v2"
	"os"
	"fossil/src/engine/util"
	"fossil/src/engine/client"
	"fmt"
	"strconv"
	"time"
	"text/tabwriter"
)

func main() {
	optUsername := getopt.StringLong("user",'u',"","Username")
	optPassword := getopt.StringLong("pass",'s',"","Password")
	optProfile := getopt.StringLong("profile",'p',"","Profile name")
	optConfig := getopt.StringLong("config",'c',"","Config name")
	optConfigPath := getopt.StringLong("configPath",'o',"","Path to configs directory")
	optConfigFile := getopt.StringLong("configFile",'f',"","Path to config file")
	optPolicy := getopt.StringLong("policy",'i',"","Backup policy as defined in config")
	optAction := getopt.StringLong("action",'a',"","backup|backupList|listProfiles|listConfigs|listPluginConfigs|addConfig|addPluginConfig|deleteConfig|addProfile|addSchedule|deleteSchedule|listSchedules|jobStatus|appPluginList|storagePluginList|archivePluginList|pluginInfo|status")
	optPluginName := getopt.StringLong("plugin",'l',"","Name of plugin")
	optPluginType := getopt.StringLong("pluginType",'t',"","Plugin type app|storage|archive")
	optWorkflowId := getopt.StringLong("workflowId",'w',"","Workflow Id")
	optCronSchedule := getopt.StringLong("cronSchedule",'r',"","Cron Schedule Format - (min) (hour) (dayOfMOnth) (month) (dayOfWeek)")
	optLocalConfig := getopt.BoolLong("local", 0,"Use a local configuration file")
	optListSchedules := getopt.BoolLong("list-schedules", 0,"List schedules")
	optGetDefaultConfig := getopt.BoolLong("get-default-config", 0,"Get the default config file")
	optGetDefaultPluginConfig := getopt.BoolLong("get-default-plugin-config", 0,"Get the default config file")
	optGetConfig := getopt.BoolLong("get-config", 0,"Get config file")
	optGetPluginConfig := getopt.BoolLong("get-plugin-config", 0,"Get plugin config file")
	optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

    if *optHelp {
        getopt.Usage()
        os.Exit(0)
	}

	if getopt.IsSet("user") != true {
		fmt.Println("ERROR: Missing parameter --user")
		getopt.Usage()
		os.Exit(1)
	}

	if getopt.IsSet("pass") != true {
		fmt.Println("ERROR: Missing parameter --pass")
		getopt.Usage()
		os.Exit(1)
	}

	var auth client.Auth
	auth.Username = *optUsername
	auth.Password = *optPassword

	if *optGetDefaultConfig {
		config,err := client.GetDefaultConfig(auth)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}

		fmt.Println("### Default Config ###")
		fmt.Println("PluginDir = " + "\""+ config.PluginDir + "\"")
		fmt.Println("AppPlugin = " + "\"" + config.AppPlugin + "\"")
		fmt.Println("PreAppQuiesceCmd = " + "\"" + config.PreAppQuiesceCmd + "\"")
		fmt.Println("PostAppQuiesceCmd = " + "\"" + config.PostAppQuiesceCmd + "\"")
		fmt.Println("BackupCreateCmd = " + "\"" + config.BackupCreateCmd + "\"")
		fmt.Println("BackupDeleteCmd = " + "\"" + config.BackupDeleteCmd + "\"")
		fmt.Println("PreAppUnquiesceCmd = " + "\"" + config.PreAppUnquiesceCmd + "\"")
		fmt.Println("AppUnquiesceCmd = " + "\"" + config.AppUnquiesceCmd + "\"")
		fmt.Println("AppUnquiesceCmd = " + "\"" + config.AppUnquiesceCmd + "\"")
		fmt.Println("PostAppUnquiesceCmd = " + "\"" + config.PostAppUnquiesceCmd + "\"")
		fmt.Println("SendTrapErrorCmd = " + "\"" +config.SendTrapErrorCmd + "\"")
		fmt.Println("SendTrapSuccessCmd = " + "\"" + config.SendTrapSuccessCmd + "\"" + "\n")

		for _, retention := range config.BackupRetentions {
			fmt.Println("[[BackupRetentions]]")
			fmt.Println("Policy = " + "\"" + retention.Policy + "\"")
			fmt.Println("RetentionDays = " + strconv.Itoa(retention.RetentionDays) + "\n")
		}
		os.Exit(0)
	}

	if *optListSchedules {
		fmt.Println("### Job Schedules ###")
		jobScheduleResult,err := client.ListSchedules(auth)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		checkResult(jobScheduleResult.Result)

		// print friendly columns
		tw := new(tabwriter.Writer)
		tw.Init(os.Stdout, 10, 20, 5, ' ', 0)
		fmt.Fprintln(tw, "CronSchedule\t ProfileName\t ConfigName\t Policy\t")
		for _,schedule := range jobScheduleResult.JobSchedules {	
			fmt.Fprintln(tw, schedule.CronSchedule + "\t",schedule.ProfileName + "\t",schedule.ConfigName + "\t",schedule.BackupPolicy + "\t")
		}		
		tw.Flush()
	
		os.Exit(0)
	}	
	
	if *optGetDefaultPluginConfig {
		if getopt.IsSet("plugin") != true {
			fmt.Println("ERROR: Missing parameter --plugin")
			getopt.Usage()
			os.Exit(1)
		}

		configMap,err := client.GetDefaultPluginConfig(auth,*optPluginName)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}

		fmt.Println("### Default Plugin Config ###")
		for k,v := range configMap {
			fmt.Println(k + " = " + "\""+ v + "\"")
		}		
		os.Exit(0)
	}

	if *optGetConfig {
		config,err := client.GetConfig(auth,string(*optProfile),string(*optConfig))
		if err != nil {
			fmt.Println(err,"\n" + "ERROR: Couldn't get config [" + string(*optProfile) + "] config [" + string(*optConfig) + "! " + err.Error())
			os.Exit(1)	
		}
		buf,err := util.EncodeConfig(config)
		if err != nil {
			fmt.Println(err,"\n" + "ERROR: Couldn't encode config [" + string(*optProfile) + "] config [" + string(*optConfig) + "! " + err.Error())
			os.Exit(1)	
		}		
		fmt.Println(buf.String())
		os.Exit(0)		
	}	

	if *optGetPluginConfig {
		if getopt.IsSet("plugin") != true {
			fmt.Println("ERROR: Missing parameter --plugin")
			getopt.Usage()
			os.Exit(1)
		}

		pluginConfigMap,err := client.GetPluginConfig(auth,string(*optProfile),string(*optConfig),string(*optPluginName))
		if err != nil {
			fmt.Println(err,"\n" + "ERROR: Couldn't get config [" + string(*optProfile) + "] config [" + string(*optPluginName) + "! " + err.Error())
			os.Exit(1)	
		}
		
		buf,err := util.EncodePluginConfig(pluginConfigMap)
		if err != nil {
			fmt.Println(err,"\n" + "ERROR: Couldn't encode config [" + string(*optProfile) + "] config [" + string(*optConfig) + "! " + err.Error())
			os.Exit(1)	
		}	
		fmt.Println(buf.String())	
		os.Exit(0)
	}		

	if getopt.IsSet("action") != true {
		fmt.Println("ERROR: missing parameter --action")
		getopt.Usage()
		os.Exit(1)
	}

	if *optAction == "listProfiles" {
		fmt.Println("### Profile List ###")
		result,err := client.ListProfiles(auth)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		checkResult(result)
		for _,profile := range result.Data {
			fmt.Println(profile)
		}
		os.Exit(0)
	}	

	if getopt.IsSet("profile") != true {
		fmt.Println("ERROR: missing parameter --profile")
		getopt.Usage()
		os.Exit(1)
	}	
	
	if *optAction == "listConfigs" {
		fmt.Println("### Config List ###")
		result,err := client.ListConfigs(auth,string(*optProfile))
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		checkResult(result)
		for _,config := range result.Data {
			fmt.Println(config)
		}
		os.Exit(0)
	}	

	if *optAction == "addProfile" {
		result,err := client.AddProfile(auth,string(*optProfile))
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		printResult(result)
		os.Exit(0)
	}	
	if *optAction == "deleteProfile" {
		result,err := client.DeleteProfile(auth,string(*optProfile))
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		printResult(result)
		os.Exit(0)
	}
	
	if getopt.IsSet("config") != true {
		fmt.Println("ERROR: missing parameter --config")
		getopt.Usage()
		os.Exit(1)
	}	

	if *optAction == "listPluginConfigs" {
		fmt.Println("### Config List ###")
		result,err := client.ListPluginConfigs(auth,string(*optProfile),string(*optConfig))
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		checkResult(result)
		for _,plugin := range result.Data {
			fmt.Println(plugin)
		}
		os.Exit(0)
	}
	
	if *optAction == "addConfig" {
		if getopt.IsSet("configFile") != true {
			fmt.Println("ERROR: Missing parameter --configFile")
			getopt.Usage()
			os.Exit(1)
		}

		config,err := util.ReadConfig(*optConfigFile)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		result,err := client.AddConfig(auth,string(*optProfile),string(*optConfig),config)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		printResult(result)		
		os.Exit(0)
	}
	if *optAction == "addPluginConfig" {
		if getopt.IsSet("configFile") != true {
			fmt.Println("ERROR: Missing parameter --configFile")
			getopt.Usage()
			os.Exit(1)
		}

		configMap,err := util.ReadConfigToMap(*optConfigFile)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		result,err := client.AddPluginConfig(auth,string(*optProfile),string(*optConfig),string(*optPluginName),configMap)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		printResult(result)		
		os.Exit(0)
	}					
	if *optAction == "deleteConfig" {
		result,err := client.DeleteConfig(auth,string(*optProfile),string(*optConfig))
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		printResult(result)	
		os.Exit(0)
	}

	// Get config
	var config util.Config
	var err error
	if *optLocalConfig {	
		var configPath string
		var configDir string
		if getopt.IsSet("configPath") == true {
			configPath = *optConfigPath + "/" + *optProfile + "/" + *optConfig + "/" + *optConfig + ".conf"
			configDir = *optConfigPath + "/" + *optProfile + "/" + *optConfig 
		} else {
			configPath = "configs/" + *optProfile + "/" + *optConfig + "/" + *optConfig + ".conf"
			configDir = "configs/" + *optProfile + "/" + *optConfig
		}

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Println(err,"\n" + "ERROR: Profile of Config don't exist")
			os.Exit(1)
		}

		//read config file into struct
		config,err = util.ReadConfig(configPath)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}

		//load dynamic plugin parameters into config struct
		if config.AppPlugin != "" {
			var err error
			appConfigPath := configDir + "/" + config.AppPlugin + ".conf"
			config,err = util.SetAppPluginParameters(appConfigPath, config)
			if err != nil {
				fmt.Println("ERROR: " + err.Error())
				os.Exit(1)
			}
		}
		if config.StoragePlugin != "" {
			var err error
			storageConfigPath := configDir + "/" + config.StoragePlugin + ".conf"
			config,err = util.SetStoragePluginParameters(storageConfigPath, config)
			if err != nil {
				fmt.Println("ERROR: " + err.Error())
				os.Exit(1)
			}
		}
	} else {
		config,err = client.GetConfig(auth,string(*optProfile),string(*optConfig))
		if err != nil {
			fmt.Println(err,"\n" + "ERROR: Couldn't get profile [" + string(*optProfile) + "] config [" + string(*optConfig) + "! " + err.Error())
			os.Exit(1)	
		}

		//load dynamic plugin parameters into config struct
		if config.AppPlugin != "" {
			appConfigMap,err := client.GetPluginConfig(auth,string(*optProfile),string(*optConfig),config.AppPlugin)
			if err != nil {
				fmt.Println(err,"\n" + "ERROR: Couldn't get profile [" + string(*optProfile) + "] config [" + config.AppPlugin + "! " + err.Error())
				os.Exit(1)	
			}	
			config.AppPluginParameters = appConfigMap		
		}
		if config.StoragePlugin != "" {
			storageConfigMap,err := client.GetPluginConfig(auth,string(*optProfile),string(*optConfig),config.StoragePlugin)
			if err != nil {
				fmt.Println(err,"\n" + "ERROR: Couldn't get profile [" + string(*optProfile) + "] config [" + config.StoragePlugin + "! " + err.Error())
				os.Exit(1)	
			}	
			config.StoragePluginParameters = storageConfigMap	
		}		
	}	

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

	fmt.Println("########## Welcome to Fossil Framework ##########")

	if *optAction == "backup" {
		logger := util.GetLoggerInstance()

		workflowResult,err := client.StartBackupWorkflow(auth,string(*optProfile),string(*optConfig),string(*optPolicy),config)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		
		util.LogResult(logger, workflowResult.Result)
		if workflowResult.Result.Code != 0 {
			os.Exit(1)
		}

		workflowId := workflowResult.Id
		var completedSteps []int
		// loop and wait for all workflow steps to complete
		for {
			time.Sleep(1 * time.Second)
			workflow,err := client.GetWorkflowStatus(auth,*optProfile,*optConfig,workflowId)
			if err != nil {
				fmt.Println("ERROR: " + err.Error())
				os.Exit(1)
			}

			// Print results for a step only once
			for _, step := range workflow.Steps {
				if step.Status == "COMPLETE" || step.Status == "ERROR" {
					if !util.IntInSlice(step.Id,completedSteps) {
						completedSteps = append(completedSteps,step.Id)
						results,err := client.GetWorkflowStepResults(auth,*optProfile,*optConfig,workflowId,step.Id)
						if err != nil {
							fmt.Println("ERROR: " + err.Error())
							os.Exit(1)
						}
						util.LogResults(logger, results)
					}
				}
			}

			if workflow.Status == "COMPLETE" || workflow.Status == "ERROR"  {
				break
			}
			time.Sleep(4 * time.Second)
		}
	} else if *optAction == "backupList" {
		msg := fmt.Sprintf("### List of Backups for policy [%s] ###",*optPolicy)
		fmt.Println(msg)

		backups,err := client.BackupList(auth,string(*optProfile),string(*optConfig),string(*optPolicy),config)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}

		backupsByPolicy := util.GetBackupsByPolicy(string(*optPolicy),backups.Backups)
		checkResult(backups.Result)

		for _, backup := range backupsByPolicy {
			fmt.Println(backup.Name, backup.Policy, backup.WorkflowId, backup.Timestamp)
		}
	} else if *optAction == "jobList" {	
		msg := fmt.Sprintf("### List of Jobs for profile [%s] config [%s] ###",*optProfile,*optConfig)
		fmt.Println(msg)

		jobs,err := client.GetJobList(auth,string(*optProfile),string(*optConfig))
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		checkResult(jobs.Result)

		// print friendly columns
		tw := new(tabwriter.Writer)
		tw.Init(os.Stdout, 10, 20, 5, ' ', 0)
		fmt.Fprintln(tw, "WorkflowId\t Status\t Start Time\t")
		for _, job := range jobs.Jobs {
			fmt.Fprintln(tw, util.IntToString(job.Id) + "\t",job.Status + "\t",job.Timestamp + "\t")
		}		
		tw.Flush()
	} else if *optAction == "jobStatus" {
		if getopt.IsSet("workflowId") != true {
			fmt.Println("ERROR: Missing parameter --workflowId")
			getopt.Usage()
			os.Exit(1)
		}

		logger := util.GetLoggerInstance()

		workflowId := util.StringToInt(*optWorkflowId)
		var completedSteps []int
		// loop and wait for all workflow steps to complete
		for {
			time.Sleep(1 * time.Second)
			workflow,err := client.GetWorkflowStatus(auth,*optProfile,*optConfig,workflowId)
			if err != nil {
				fmt.Println("ERROR: " + err.Error())
				os.Exit(1)
			}

			// Print results for a step only once
			for _, step := range workflow.Steps {
				if step.Status == "COMPLETE" || step.Status == "ERROR" {
					if !util.IntInSlice(step.Id,completedSteps) {
						completedSteps = append(completedSteps,step.Id)
						results,err := client.GetWorkflowStepResults(auth,*optProfile,*optConfig,workflowId,step.Id)
						if err != nil {
							fmt.Println("ERROR: " + err.Error())
							os.Exit(1)
						}

						util.LogResults(logger, results)
					}
				}
			}

			if workflow.Status == "COMPLETE" || workflow.Status == "ERROR"  {
				break
			}
			time.Sleep(4 * time.Second)
		}			
	} else if *optAction == "appPluginList" {
		fmt.Println("### List of Application Plugins ###")

		var plugins []string
		appPlugins,err := client.AppPluginList(auth,"app",config)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}

		plugins = util.JoinArray(appPlugins,plugins)

		for _, plugin := range plugins {
			fmt.Println(plugin)
		}
	} else if *optAction == "storagePluginList" {
		fmt.Println("### List of Storage Plugins ###")
	
		var plugins []string
		storagePlugins,err := client.StoragePluginList(auth,"storage",config)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}

		plugins = util.JoinArray(storagePlugins,plugins)
	
		for _, plugin := range plugins {
			fmt.Println(plugin)
		}
	} else if *optAction == "archivePluginList" {
		fmt.Println("### List of Archive Plugins ###")
		
		var plugins []string
		archivePlugins,err := client.ArchivePluginList(auth,"archive",config)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}

		plugins = util.JoinArray(archivePlugins,plugins)
		
		for _, plugin := range plugins {
			fmt.Println(plugin)
		}					
	} else if *optAction == "pluginInfo" {
		if getopt.IsSet("plugin") != true {
			fmt.Println("ERROR: Missing parameter --plugin")
			getopt.Usage()
			os.Exit(1)
		}

		if getopt.IsSet("pluginType") != true {
			fmt.Println("ERROR: Missing parameter --pluginType")
			getopt.Usage()
			os.Exit(1)
		}

		var pluginName string = *optPluginName
		var pluginInfoResult util.PluginInfoResult
		var err error

		if *optPluginType == "app" {
			pluginInfoResult,err = client.AppPluginInfo(auth,config,pluginName,*optPluginType)
			if err != nil {
				fmt.Println("ERROR: " + err.Error())
				os.Exit(1)
			}
		} else if *optPluginType == "storage" {
			pluginInfoResult,err = client.StoragePluginInfo(auth,config,pluginName,*optPluginType)
			if err != nil {
				fmt.Println("ERROR: " + err.Error())
				os.Exit(1)
			}
		} else if *optPluginType == "archive" {
			pluginInfoResult,err = client.ArchivePluginInfo(auth,config,pluginName,*optPluginType)
			if err != nil {
				fmt.Println("ERROR: " + err.Error())
				os.Exit(1)
			}
		} else {
			error := fmt.Sprintf("ERROR: Plugin type must be app|storage|archive")
			fmt.Println(error)
		}	

		checkResult(pluginInfoResult.Result)

		fmt.Println("### Plugin Information ###")
		fmt.Println("Name:", pluginInfoResult.Plugin.Name)
		fmt.Println("Description:", pluginInfoResult.Plugin.Description)
		fmt.Println("Type:", pluginInfoResult.Plugin.Type)
		fmt.Println("Capabilities:", pluginInfoResult.Plugin.Capabilities)
	} else if *optAction == "status" {
		fmt.Println("### Checking status of services ###")

		workflowStatus,err := client.GetWorkflowServiceStatus(auth)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		fmt.Println("Workflow Service:", workflowStatus)

		appStatus,err := client.GetAppServiceStatus(auth)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		fmt.Println("App Service:", appStatus)

		storageStatus,err := client.GetStorageServiceStatus(auth)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		fmt.Println("Storage Service:", storageStatus)
	} else if *optAction == "addSchedule" {
		if getopt.IsSet("cronSchedule") != true {
			fmt.Println("ERROR: Missing parameter --cronSchedule")
			getopt.Usage()
			os.Exit(1)
		}

		result,err := client.AddSchedule(auth,*optProfile,*optConfig,*optPolicy,*optCronSchedule)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		printResult(result)
		os.Exit(0)
	} else if *optAction == "deleteSchedule" {
		result,err := client.DeleteSchedule(auth,*optProfile,*optConfig,*optPolicy)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
		printResult(result)
		os.Exit(0)	
	} else {
		fmt.Println("ERROR: incorrect parameter", *optAction)
		getopt.Usage()
		os.Exit(1)
	}
}

func checkResult(result util.Result) {
	logger := util.GetLoggerInstance()
	if result.Code != 0 {
		util.LogResult(logger, result)
		os.Exit(1)
	}
}

func printResult(result util.Result) {
	logger := util.GetLoggerInstance()
	util.LogResult(logger, result)
}