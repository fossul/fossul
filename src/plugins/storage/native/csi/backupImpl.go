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
	"fossul/src/client/k8s"
	"fossul/src/engine/util"
)

func (s storagePlugin) Backup(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	timestampToString := fmt.Sprintf("%d", config.WorkflowTimestamp)
	//backupName := util.GetBackupName(config.StoragePluginParameters["PvcName"], config.SelectedBackupPolicy, config.WorkflowId, timestampToString)
	backupName := fmt.Sprintf(config.StoragePluginParameters["PvcName"] + "-" + config.SelectedBackupPolicy + "-" + config.WorkflowId + "-" + timestampToString)
	timeout := util.StringToInt(config.StoragePluginParameters["SnapshotTimeoutSeconds"])

	msg := util.SetMessage("INFO", "Creating CSI snapshot ["+backupName+"] of pvc ["+config.StoragePluginParameters["PvcName"]+"] namespace ["+config.StoragePluginParameters["Namespace"]+"] snapshot class ["+config.StoragePluginParameters["SnapshotClass"]+" ] timeout ["+config.StoragePluginParameters["SnapshotTimeoutSeconds"]+"]")
	messages = append(messages, msg)

	err := k8s.CreateSnapshot(backupName, config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["SnapshotClass"], config.StoragePluginParameters["PvcName"], config.AccessWithinCluster, timeout)

	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	msg = util.SetMessage("INFO", "CSI snapshot created successfully")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}
