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
	"github.com/fossul/fossul/src/client/k8s"
	"github.com/fossul/fossul/src/engine/util"
	"github.com/fossul/fossul/src/plugins/pluginUtil"
)

func (s storagePlugin) BackupList(config util.Config) util.Backups {
	var backups util.Backups
	var result util.Result
	var messages []util.Message

	snapshots, err := k8s.ListSnapshots(config.StoragePluginParameters["Namespace"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)
		backups.Result = result

		return backups
	}

	var snapshotList []string
	for _, snapshot := range snapshots.Items {
		snapshotList = append(snapshotList, snapshot.Name)
	}

	backupList, err := pluginUtil.ListSnapshots(snapshotList, config.StoragePluginParameters["BackupName"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)
		backups.Result = result

		return backups
	}

	result = util.SetResult(0, messages)
	backups.Result = result
	backups.Backups = backupList

	return backups
}
