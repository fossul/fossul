package util

import (
	"testing"
	"log"
)

func TestReadConfig(t *testing.T) {
	configFile := "../../cli/configs/default/default.conf"

	config,err := ReadConfig(configFile)
	if err != nil {
		t.Fail()
	}

	log.Println("Config Struct",config)

	if config.AppPlugin != "sample-app" {
		t.Fail()
	}

	if len(config.BackupRetentions) != 2 {
		t.Fail()
	}
}

func TestReadConfigToMap(t *testing.T) {
	configFile := "../../cli/configs/default/default.conf"

	configMap,err := ReadConfigToMap(configFile)
	if err != nil {
		t.Fail()
	}

	log.Println("Config Map",configMap)

	if configMap["AppPlugin"] != "sample-app" {
		t.Fail()
	}

	if configMap["StoragePlugin"] != "sample-storage" {
		t.Fail()
	}
}	

func TestSetAppPluginParameters(t *testing.T) {
	configFile := "../../cli/configs/default/default.conf"
	appPluginFile := "../../cli/configs/default/sample-app.conf"

	config,err := ReadConfig(configFile)
	if err != nil {
		t.Fail()
	}

	config,err = SetAppPluginParameters(appPluginFile,config)
	if err != nil {
		t.Fail()
	}

	if config.AppPluginParameters["SampleAppVar1"] != "foo" {
		t.Fail()
	}

	if config.AppPluginParameters["SampleAppVar2"] != "bar" {
		t.Fail()
	}
}	

func TestSetStoragePluginParameters(t *testing.T) {
	configFile := "../../cli/configs/default/default.conf"
	storagePluginFile := "../../cli/configs/default/sample-storage.conf"

	config,err := ReadConfig(configFile)
	if err != nil {
		t.Fail()
	}

	config,err = SetStoragePluginParameters(storagePluginFile,config)
	if err != nil {
		t.Fail()
	}

	if config.StoragePluginParameters["SampleStorageVar1"] != "foo" {
		t.Fail()
	}

	if config.StoragePluginParameters["SampleStorageVar2"] != "bar" {
		t.Fail()
	}
}	

func TestExistsBackupRetention(t *testing.T) {
	configFile := "../../cli/configs/default/default.conf"
	config,err := ReadConfig(configFile)
	if err != nil {
		t.Fail()
	}

	exists := ExistsBackupRetention("daily",config.BackupRetentions)

	if !exists {
		t.Fail()
	}
}	

func TestGetBackupRetention(t *testing.T) {
	configFile := "../../cli/configs/default/default.conf"
	config,err := ReadConfig(configFile)
	if err != nil {
		t.Fail()
	}

	retentionDaily := GetBackupRetention("daily",config.BackupRetentions)

	retentionWeekly := GetBackupRetention("weekly",config.BackupRetentions)

	if retentionDaily != 5 {
		t.Fail()
	}

	if retentionWeekly != 2 {
		t.Fail()
	}
}	

