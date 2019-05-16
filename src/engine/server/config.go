package main

import (
	"fossil/src/engine/util"
)

func GetConsolidatedConfig(profileName,configName,policyName string) (util.Config,error) {
	conf := configDir + "/" + profileName + "/" + configName + "/" + configName + ".conf"
	config,err := util.ReadConfig(conf)
	
	if err != nil {
		return config,err
	}

	config.ProfileName = profileName
	config.ConfigName = configName

	backupRetention := util.GetBackupRetention(policyName,config.BackupRetentions)
	config.SelectedBackupRetention = backupRetention
	config.SelectedBackupPolicy = policyName

	if config.AppPlugin != "" {
		appConf := configDir + "/" + profileName + "/" + configName + "/" + config.AppPlugin + ".conf"
		appConfigMap,err := util.ReadConfigToMap(appConf)

		if err != nil {
			return config,err
		}

		config.AppPluginParameters = appConfigMap
	}
	
	if config.StoragePlugin != "" {
		storageConf := configDir + "/" + profileName + "/" + configName + "/" + config.StoragePlugin + ".conf"
		storageConfigMap,err := util.ReadConfigToMap(storageConf)

		if err != nil {
			return config,err
		}

		config.StoragePluginParameters = storageConfigMap
	}	

	return config,nil
}