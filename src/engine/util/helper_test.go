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

	if !strings.Contains(backupPath, "/backupdest/default/default/mybackup_daily_777") {
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
