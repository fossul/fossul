package main

import (
	"os"
	"fossil/src/engine/util"
	"fossil/src/engine/client"
	"fmt"
	"strconv"
	"time"
	"text/tabwriter"
	"errors"
)

func WriteCredentialFile(credentialFile,serverHostname,serverPort,appHostname,appPort,storageHostname,storagePort,username,password string) {
	var auth client.Auth
	auth.ServerHostname = serverHostname
	auth.ServerPort = serverPort
	auth.AppHostname = appHostname
	auth.AppPort = appPort
	auth.StorageHostname = storageHostname
	auth.StoragePort = storagePort
	auth.Username = username
	auth.Password = password

	err := util.WriteGob(credentialFile,auth)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func ReadCredentialFile(credentialFile string) (client.Auth) {
	auth := &client.Auth{}

	err := util.ReadGob(credentialFile,&auth)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	return *auth
}

func GetDefaultConfig(auth client.Auth) {
	config,err := client.GetDefaultConfig(auth)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("### Default Config ###")
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

func ListSchedules(auth client.Auth) {
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

func GetDefaultPluginConfig(auth client.Auth,pluginName string) {

	configMap,err := client.GetDefaultPluginConfig(auth,pluginName)
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

func GetConfig(auth client.Auth,profileName,configName string) {
	config,err := client.GetConfig(auth,profileName,configName)
	if err != nil {
		fmt.Println(err,"\n" + "ERROR: Couldn't get config [" + profileName + "] config [" + configName + "! " + err.Error())
		os.Exit(1)	
	}
	buf,err := util.EncodeConfig(config)
	if err != nil {
		fmt.Println(err,"\n" + "ERROR: Couldn't encode config [" + profileName + "] config [" + configName + "! " + err.Error())
		os.Exit(1)	
	}		
	fmt.Println(buf.String())
	os.Exit(0)		
}	

func GetPluginConfig(auth client.Auth,profileName,configName,pluginName string) {
	pluginConfigMap,err := client.GetPluginConfig(auth,profileName,configName,pluginName)
	if err != nil {
		fmt.Println(err,"\n" + "ERROR: Couldn't get config [" + profileName + "] config [" + configName + "! " + err.Error())
		os.Exit(1)	
	}
		
	buf,err := util.EncodePluginConfig(pluginConfigMap)
	if err != nil {
		fmt.Println(err,"\n" + "ERROR: Couldn't encode config [" + profileName + "] config [" + configName + "! " + err.Error())
		os.Exit(1)	
	}	
	fmt.Println(buf.String())	
	os.Exit(0)
}	

func ListProfiles(auth client.Auth) {
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

func ListConfigs(auth client.Auth,profileName string) {
	fmt.Println("### Config List ###")
	result,err := client.ListConfigs(auth,profileName)
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

func AddProfile(auth client.Auth,profileName string) {
	result,err := client.AddProfile(auth,profileName)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
	printResult(result)
	os.Exit(0)
}	

func DeleteProfile(auth client.Auth,profileName string) {
	result,err := client.DeleteProfile(auth,profileName)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
	printResult(result)
	os.Exit(0)
}
	

func ListPluginConfigs(auth client.Auth,profileName,configName string) {
	fmt.Println("### Config List ###")
	result,err := client.ListPluginConfigs(auth,profileName,configName)
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

func AddConfig(auth client.Auth,profileName,configName,configFile string) {
	config,err := util.ReadConfig(configFile)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
	result,err := client.AddConfig(auth,profileName,configName,config)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
	printResult(result)		
	os.Exit(0)
}

func AddPluginConfig(auth client.Auth,profileName,configName,pluginName,configFile string) {

	configMap,err := util.ReadConfigToMap(configFile)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
	result,err := client.AddPluginConfig(auth,profileName,configName,pluginName,configMap)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
	printResult(result)		
	os.Exit(0)
}

func DeleteConfig(auth client.Auth,profileName,configName string) {
	result,err := client.DeleteConfig(auth,profileName,configName)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
	printResult(result)	
	os.Exit(0)
}

func ImportLocalConfig(configDir,configPath string) (util.Config,error) {
	config,err := util.ReadConfig(configPath)
	if err != nil {
		return config,err
	}

	//load dynamic plugin parameters into config struct
	if config.AppPlugin != "" {
		var err error
		appConfigPath := configDir + "/" + config.AppPlugin + ".conf"
		config,err = util.SetAppPluginParameters(appConfigPath, config)
		if err != nil {
			return config,err
		}
	}
	if config.StoragePlugin != "" {
		var err error
		storageConfigPath := configDir + "/" + config.StoragePlugin + ".conf"
		config,err = util.SetStoragePluginParameters(storageConfigPath, config)
		if err != nil {
			return config,err
		}
	}

	return config,nil
}

func ImportServerConfig(auth client.Auth,profileName,configName string) (util.Config,error) {
	config,err := client.GetConfig(auth,profileName,configName)
	if err != nil {
		return config,errors.New("ERROR: Couldn't get profile [" + profileName + "] config [" + configName + "! " + err.Error())
	}

	//load dynamic plugin parameters into config struct
	if config.AppPlugin != "" {
		appConfigMap,err := client.GetPluginConfig(auth,profileName,configName,config.AppPlugin)
		if err != nil {
			return config,errors.New("ERROR: Couldn't get profile [" + profileName + "] config [" + config.AppPlugin + "! " + err.Error())
		}	
		config.AppPluginParameters = appConfigMap		
	}
	if config.StoragePlugin != "" {
		storageConfigMap,err := client.GetPluginConfig(auth,profileName,configName,config.StoragePlugin)
		if err != nil {
			return config,errors.New("ERROR: Couldn't get profile [" + profileName + "] config [" + config.StoragePlugin + "! " + err.Error())
		}	
		config.StoragePluginParameters = storageConfigMap	
	}
	
	return config,nil
}	

func BackupWithLocalConfig(auth client.Auth,profileName,configName,policyName string,config util.Config) {
	logger := util.GetLoggerInstance()

	workflowResult,err := client.StartBackupWorkflowLocalConfig(auth,profileName,configName,policyName,config)
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
		workflow,err := client.GetWorkflowStatus(auth,profileName,configName,workflowId)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}

		// Print results for a step only once
		for _, step := range workflow.Steps {
			if step.Status == "COMPLETE" || step.Status == "ERROR" {
				if !util.IntInSlice(step.Id,completedSteps) {
					completedSteps = append(completedSteps,step.Id)
					results,err := client.GetWorkflowStepResults(auth,profileName,configName,workflowId,step.Id)
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
}

func Backup(auth client.Auth,profileName,configName,policyName string) {
	logger := util.GetLoggerInstance()

	workflowResult,err := client.StartBackupWorkflow(auth,profileName,configName,policyName)
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
		workflow,err := client.GetWorkflowStatus(auth,profileName,configName,workflowId)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}

		// Print results for a step only once
		for _, step := range workflow.Steps {
			if step.Status == "COMPLETE" || step.Status == "ERROR" {
				if !util.IntInSlice(step.Id,completedSteps) {
					completedSteps = append(completedSteps,step.Id)
					results,err := client.GetWorkflowStepResults(auth,profileName,configName,workflowId,step.Id)
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
}

func BackupList(auth client.Auth,profileName,configName,policyName string,config util.Config) {
	msg := fmt.Sprintf("### List of Backups for policy [%s] ###",policyName)
	fmt.Println(msg)

	backups,err := client.BackupList(auth,profileName,configName,policyName,config)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	backupsByPolicy := util.GetBackupsByPolicy(policyName,backups.Backups)
	checkResult(backups.Result)

	for _, backup := range backupsByPolicy {
		fmt.Println(backup.Name, backup.Policy, backup.WorkflowId, backup.Timestamp)
	}
}	

func JobList(auth client.Auth,profileName,configName string) {
	msg := fmt.Sprintf("### List of Jobs for profile [%s] config [%s] ###",profileName,configName)
	fmt.Println(msg)

	jobs,err := client.GetJobList(auth,profileName,configName)
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
}

func JobStatus(auth client.Auth,profileName,configName,workflowId string) {
	logger := util.GetLoggerInstance()

	workflowIdInt := util.StringToInt(workflowId)
	var completedSteps []int
	// loop and wait for all workflow steps to complete
	for {
		time.Sleep(1 * time.Second)
		workflow,err := client.GetWorkflowStatus(auth,profileName,configName,workflowIdInt)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}

		// Print results for a step only once
		for _, step := range workflow.Steps {
			if step.Status == "COMPLETE" || step.Status == "ERROR" {
				if !util.IntInSlice(step.Id,completedSteps) {
					completedSteps = append(completedSteps,step.Id)
					results,err := client.GetWorkflowStepResults(auth,profileName,configName,workflowIdInt,step.Id)
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
}

func AppPluginList(auth client.Auth) {
	fmt.Println("### List of Application Plugins ###")

	var plugins []string
	appPlugins,err := client.AppPluginList(auth,"app")
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	plugins = util.JoinArray(appPlugins,plugins)

	for _, plugin := range plugins {
		fmt.Println(plugin)
	}
	os.Exit(0)
}

func StoragePluginList(auth client.Auth) {
	fmt.Println("### List of Storage Plugins ###")
	
	var plugins []string
	storagePlugins,err := client.StoragePluginList(auth,"storage")
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	plugins = util.JoinArray(storagePlugins,plugins)
	
	for _, plugin := range plugins {
		fmt.Println(plugin)
	}
	os.Exit(0)
}

func ArchivePluginList(auth client.Auth) {
	fmt.Println("### List of Archive Plugins ###")
		
	var plugins []string
	archivePlugins,err := client.ArchivePluginList(auth,"archive")
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	plugins = util.JoinArray(archivePlugins,plugins)
		
	for _, plugin := range plugins {
		fmt.Println(plugin)
	}
	os.Exit(0)
}

func PluginInfo(auth client.Auth,config util.Config,pluginName,pluginType string) {
	var pluginInfoResult util.PluginInfoResult
	var err error

	if pluginType == "app" {
		pluginInfoResult,err = client.AppPluginInfo(auth,config,pluginName,pluginType)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
	} else if pluginType == "storage" {
		pluginInfoResult,err = client.StoragePluginInfo(auth,config,pluginName,pluginType)
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			os.Exit(1)
		}
	} else if pluginType == "archive" {
		pluginInfoResult,err = client.ArchivePluginInfo(auth,config,pluginName,pluginType)
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
}

func Status(auth client.Auth) {
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

	os.Exit(0)
}

func AddSchedule(auth client.Auth,profileName,configName,policyName,cronSchedule string) {
	result,err := client.AddSchedule(auth,profileName,configName,policyName,cronSchedule)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
	printResult(result)
	os.Exit(0)
}

func DeleteSchedule(auth client.Auth,profileName,configName,policyName string) {
	result,err := client.DeleteSchedule(auth,profileName,configName,policyName)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
	printResult(result)
	os.Exit(0)
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