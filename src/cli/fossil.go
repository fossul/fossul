package main

import (
	"github.com/pborman/getopt/v2"
	"os"
	"fossil/src/engine/util"
	"fossil/src/engine/client"
	"fmt"
)

func main() {
	optUsername := getopt.StringLong("user",'u',"","Username")
	optPassword := getopt.StringLong("pass",'s',"","Password")
	optProfile := getopt.StringLong("profile",'p',"","Profile name")
	optConfig := getopt.StringLong("config",'c',"","Config name")
	optConfigPath := getopt.StringLong("config-path",'o',"","Path to configs directory")
	optConfigFile := getopt.StringLong("config-file",'f',"","Path to config file")
	optPolicy := getopt.StringLong("policy",'i',"","Backup policy as defined in config")
	optAction := getopt.StringLong("action",'a',"","backup|backupList|listProfiles|listConfigs|listPluginConfigs|addConfig|addPluginConfig|deleteConfig|addProfile|addSchedule|deleteSchedule|jobStatus|pluginInfo|status")
	optPluginName := getopt.StringLong("plugin",'l',"","Name of plugin")
	optPluginType := getopt.StringLong("plugin-type",'t',"","Plugin type app|storage|archive")
	optWorkflowId := getopt.StringLong("workflow-id",'w',"","Workflow Id")
	optCronSchedule := getopt.StringLong("cron-schedule",'r',"","Cron Schedule Format - (min) (hour) (dayOfMOnth) (month) (dayOfWeek)")
	optLocalConfig := getopt.BoolLong("local", 0,"Use a local configuration file")
	optListSchedules := getopt.BoolLong("list-schedules", 0,"List schedules")
	optAppPluginList := getopt.BoolLong("list-app-plugins", 0,"List app plugins")
	optStoragePluginList := getopt.BoolLong("list-storage-plugins", 0,"List storage plugins")
	optArchivePluginList := getopt.BoolLong("list-archive-plugins", 0,"List archive plugins")

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

	if *optAction == "status" {
		Status(auth)
	}

	if *optGetDefaultConfig {
		GetDefaultConfig(auth)
	}

	if *optListSchedules {
		ListSchedules(auth)
	}	
	
	if *optGetDefaultPluginConfig {
		if getopt.IsSet("plugin") != true {
			fmt.Println("ERROR: Missing parameter --plugin")
			getopt.Usage()
			os.Exit(1)
		}

		GetDefaultPluginConfig(auth,*optPluginName)
	}

	if *optGetConfig {
		if getopt.IsSet("profile") != true {
			fmt.Println("ERROR: missing parameter --profile")
			getopt.Usage()
			os.Exit(1)
		}	

		if getopt.IsSet("config") != true {
			fmt.Println("ERROR: missing parameter --config")
			getopt.Usage()
			os.Exit(1)
		}	

		GetConfig(auth,string(*optProfile),string(*optConfig))	
	}	

	if *optGetPluginConfig {
		if getopt.IsSet("profile") != true {
			fmt.Println("ERROR: missing parameter --profile")
			getopt.Usage()
			os.Exit(1)
		}	

		if getopt.IsSet("config") != true {
			fmt.Println("ERROR: missing parameter --config")
			getopt.Usage()
			os.Exit(1)
		}	
		
		if getopt.IsSet("plugin") != true {
			fmt.Println("ERROR: Missing parameter --plugin")
			getopt.Usage()
			os.Exit(1)
		}

		GetPluginConfig(auth,string(*optProfile),string(*optConfig),string(*optPluginName))
	}
	
	if *optAppPluginList {
		fmt.Println("HERE")
		AppPluginList(auth)
	}

	if *optStoragePluginList {
		StoragePluginList(auth)
	}	
	
	if *optArchivePluginList {
		ArchivePluginList(auth)		
	}

	if getopt.IsSet("action") != true {
		fmt.Println("ERROR: missing parameter --action")
		getopt.Usage()
		os.Exit(1)
	}

	if *optAction == "listProfiles" {
		ListProfiles(auth)
	}	

	if getopt.IsSet("profile") != true {
		fmt.Println("ERROR: missing parameter --profile")
		getopt.Usage()
		os.Exit(1)
	}	
	
	if *optAction == "listConfigs" {
		ListConfigs(auth,string(*optProfile))
	}	

	if *optAction == "addProfile" {
		AddProfile(auth,string(*optProfile))
	}	

	if *optAction == "deleteProfile" {
		DeleteProfile(auth,string(*optProfile))
	}
	
	if getopt.IsSet("config") != true {
		fmt.Println("ERROR: missing parameter --config")
		getopt.Usage()
		os.Exit(1)
	}	

	if *optAction == "listPluginConfigs" {
		ListPluginConfigs(auth,string(*optProfile),string(*optConfig))
	}
	
	if *optAction == "addConfig" {
		if getopt.IsSet("config-file") != true {
			fmt.Println("ERROR: Missing parameter --config-file")
			getopt.Usage()
			os.Exit(1)
		}

		AddConfig(auth,string(*optProfile),string(*optConfig),string(*optConfigFile))
	}

	if *optAction == "addPluginConfig" {
		if getopt.IsSet("plugin") != true {
			fmt.Println("ERROR: Missing parameter --plugin")
			getopt.Usage()
			os.Exit(1)
		}	

		if getopt.IsSet("config-file") != true {
			fmt.Println("ERROR: Missing parameter --config-file")
			getopt.Usage()
			os.Exit(1)
		}

		AddPluginConfig(auth,string(*optProfile),string(*optConfig),string(*optPluginName),string(*optConfigFile))
	}	

	if *optAction == "deleteConfig" {
		DeleteConfig(auth,string(*optProfile),string(*optConfig))
	}

	// Get config
	var config util.Config
	var err error
	if *optLocalConfig {	
		var configPath string
		var configDir string
		if getopt.IsSet("config-path") == true {
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

		config,err = ImportLocalConfig(configDir,configPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		config,err = ImportServerConfig(auth,string(*optProfile),string(*optConfig))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
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
			fmt.Println("ERROR: policy [" + *optPolicy + "] does not match policy defined in config")
			os.Exit(1)
		}	
	}

	fmt.Println("########## Welcome to Fossil Framework ##########")

	if *optAction == "backup" {
		if *optLocalConfig {
			BackupWithLocalConfig(auth,string(*optProfile),string(*optConfig),string(*optPolicy),config)
		} else {
			Backup(auth,string(*optProfile),string(*optConfig),string(*optPolicy))	
		}	
	} else if *optAction == "backupList" {
		BackupList(auth,string(*optProfile),string(*optConfig),string(*optPolicy),config)
	} else if *optAction == "jobList" {	
		JobList(auth,string(*optProfile),string(*optConfig))
	} else if *optAction == "jobStatus" {
		if getopt.IsSet("workflow-id") != true {
			fmt.Println("ERROR: Missing parameter --workflow-id")
			getopt.Usage()
			os.Exit(1)
		}

		JobStatus(auth,*optProfile,*optConfig,*optWorkflowId)			
	} else if *optAction == "pluginInfo" {
		if getopt.IsSet("plugin") != true {
			fmt.Println("ERROR: Missing parameter --plugin")
			getopt.Usage()
			os.Exit(1)
		}

		if getopt.IsSet("plugin-type") != true {
			fmt.Println("ERROR: Missing parameter --plugin-type")
			getopt.Usage()
			os.Exit(1)
		}

		PluginInfo(auth,config,*optPluginName,*optPluginType)
	} else if *optAction == "addSchedule" {
		if getopt.IsSet("policy") != true {
			fmt.Println("ERROR: missing parameter --policy")
			getopt.Usage()
			os.Exit(1)	
		}

		if getopt.IsSet("cron-schedule") != true {
			fmt.Println("ERROR: Missing parameter --cron-schedule")
			getopt.Usage()
			os.Exit(1)
		}

		AddSchedule(auth,*optProfile,*optConfig,*optPolicy,*optCronSchedule)
	} else if *optAction == "deleteSchedule" {
		DeleteSchedule(auth,*optProfile,*optConfig,*optPolicy)
	} else {
		fmt.Println("ERROR: incorrect parameter", *optAction)
		getopt.Usage()
		os.Exit(1)
	}
}