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
	"strings"
	"testing"
)

func TestGetBackupDir(t *testing.T) {
	configMap := getConfigMap()
	backupDir := GetBackupDirFromMap(configMap)

	if backupDir != "/backupdest/default/default" {
		t.Fail()
	}
}

func TestGetBackupPath(t *testing.T) {
	configMap := getConfigMap()
	backupPath := GetBackupPathFromMap(configMap)

	if !strings.Contains(backupPath, "/backupdest/default/default/mybackup-daily-777") {
		t.Fail()
	}
}

func getConfigMap() map[string]string {
	configMap := make(map[string]string)

	configMap["ProfileName"] = "default"
	configMap["ConfigName"] = "default"
	configMap["BackupDestPath"] = "/backupdest"
	configMap["BackupName"] = "mybackup"
	configMap["BackupPolicy"] = "daily"
	configMap["WorkflowId"] = "777"

	return configMap
}

func TestCreateDeleteDir(t *testing.T) {
	dir := "/tmp/foobar123"

	err := CreateDir(dir, 0755)
	if err != nil {
		t.Fail()
	}

	exists := ExistsPath(dir)

	if exists == true {
		err := RecursiveDirDelete(dir)
		if err != nil {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}
