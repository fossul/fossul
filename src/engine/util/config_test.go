/*
Copyright 2019 The Fossul Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package util

import (
	"log"
	"testing"
)

func TestReadConfig(t *testing.T) {
	configFile := "../../cli/configs/default/default.conf"

	config, err := ReadConfig(configFile)
	if err != nil {
		t.Fail()
	}

	log.Println("Config Struct", config)

	if config.AppPlugin != "sample-app" {
		t.Fail()
	}

	if len(config.BackupRetentions) != 2 {
		t.Fail()
	}
}

func TestReadConfigToMap(t *testing.T) {
	configFile := "../../cli/configs/default/default.conf"

	configMap, err := ReadConfigToMap(configFile)
	if err != nil {
		t.Fail()
	}

	log.Println("Config Map", configMap)

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

	config, err := ReadConfig(configFile)
	if err != nil {
		t.Fail()
	}

	config, err = SetAppPluginParameters(appPluginFile, config)
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

	config, err := ReadConfig(configFile)
	if err != nil {
		t.Fail()
	}

	config, err = SetStoragePluginParameters(storagePluginFile, config)
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

func TestSetArchivePluginParameters(t *testing.T) {
	configFile := "../../cli/configs/default/default.conf"
	archivePluginFile := "../../cli/configs/default/sample-archive.conf"

	config, err := ReadConfig(configFile)
	if err != nil {
		t.Fail()
	}

	config, err = SetArchivePluginParameters(archivePluginFile, config)
	if err != nil {
		t.Fail()
	}

	if config.ArchivePluginParameters["SampleArchiveVar1"] != "foo" {
		t.Fail()
	}

	if config.ArchivePluginParameters["SampleArchiveVar2"] != "bar" {
		t.Fail()
	}
}

func TestExistsBackupRetention(t *testing.T) {
	configFile := "../../cli/configs/default/default.conf"
	config, err := ReadConfig(configFile)
	if err != nil {
		t.Fail()
	}

	exists := ExistsBackupRetention("daily", config.BackupRetentions)

	if !exists {
		t.Fail()
	}
}

func TestGetBackupRetention(t *testing.T) {
	configFile := "../../cli/configs/default/default.conf"
	config, err := ReadConfig(configFile)
	if err != nil {
		t.Fail()
	}

	retentionDaily := GetBackupRetention("daily", config.BackupRetentions)

	retentionWeekly := GetBackupRetention("weekly", config.BackupRetentions)

	if retentionDaily != 5 {
		t.Fail()
	}

	if retentionWeekly != 2 {
		t.Fail()
	}
}
