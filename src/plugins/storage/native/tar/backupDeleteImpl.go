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
package main

import (
	"fmt"
	"os"

	"github.com/fossul/fossul/src/engine/util"
	"github.com/fossul/fossul/src/plugins/pluginUtil"
)

func (s storagePlugin) BackupDeleteWorkflow(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	backupPath := "metadata/backups/" + config.ProfileName + "/" + config.ConfigName
	backups, err := ListBackups(backupPath)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)
		return result
	}

	backupsByPolicy := util.GetBackupsByPolicy(config.SelectedBackupPolicy, backups)
	backupCount := len(backupsByPolicy)

	if backupCount > config.SelectedBackupRetention {
		count := 1
		for backup := range pluginUtil.ReverseBackupList(backupsByPolicy) {
			if count > config.SelectedBackupRetention {
				msg := util.SetMessage("INFO", fmt.Sprintf("Number of backups [%d] greater than backup retention [%d]", backupCount, config.SelectedBackupRetention))
				messages = append(messages, msg)
				backupCount = backupCount - 1

				backupName := backup.Name + "-" + backup.Policy + "-" + backup.WorkflowId + "-" + util.IntToString(backup.Epoch)
				msg = util.SetMessage("INFO", "Deleting backup "+backupName)
				messages = append(messages, msg)

				backupFile := backupPath + "/" + backupName + ".tar.gz"
				err := os.Remove(backupFile)
				if err != nil {
					msg := util.SetMessage("ERROR", "Backup "+backupName+" delete failed! "+err.Error())
					messages = append(messages, msg)
					result = util.SetResult(1, messages)
					return result
				}
				msg = util.SetMessage("INFO", "Backup "+backupName+" deleted successfully")
				messages = append(messages, msg)
			}
			count = count + 1
		}
	} else {
		msg := util.SetMessage("INFO", fmt.Sprintf("Backup deletion skipped, there are [%d] backups but backup retention is [%d]", backupCount, config.SelectedBackupRetention))
		messages = append(messages, msg)
	}

	result = util.SetResult(resultCode, messages)
	return result
}
