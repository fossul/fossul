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
	//"github.com/fossul/fossul/src/plugins/pluginUtil"
)

func (s storagePlugin) Restore(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "Performing CSI snapshot restore")
	messages = append(messages, msg)

	snapshots, err := k8s.ListSnapshots(config.StoragePluginParameters["Namespace"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	var snapshotList []string
	for _, snapshot := range snapshots.Items {
		snapshotList = append(snapshotList, snapshot.Name)
	}

	restoreSnapshot := util.GetRestoreSnapshot(config, snapshotList)
	pvcRestoreName := config.StoragePluginParameters["PvcName"] + "-restore"

	msg = util.SetMessage("INFO", "Restoring snapshot ["+restoreSnapshot+"] to new pvc ["+pvcRestoreName+"] in namespace ["+config.StoragePluginParameters["Namespace"]+"] using storage class ["+config.StoragePluginParameters["StorageClass"]+"]")
	messages = append(messages, msg)

	err = k8s.CreatePersistentVolumeClaimFromSnapshot(pvcRestoreName, config.StoragePluginParameters["PvcSize"], restoreSnapshot, config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["StorageClass"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	msg = util.SetMessage("INFO", "Updating deployment type ["+config.StoragePluginParameters["DeploymentType"]+"] deployment name  ["+config.StoragePluginParameters["DeploymentName"]+"] to use restore pvc ["+pvcRestoreName+"]")
	messages = append(messages, msg)

	if config.StoragePluginParameters["DeploymentType"] == "DeploymentConfig" {
		err := k8s.UpdateDeploymentConfigVolume(pvcRestoreName, config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["DeploymentName"], config.AccessWithinCluster)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}
	} else if config.StoragePluginParameters["DeploymentType"] == "Deployment" {
		err := k8s.UpdateDeploymentVolume(pvcRestoreName, config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["DeploymentName"], config.AccessWithinCluster)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}
	} else {
		msg := util.SetMessage("ERROR", "Couldn't find Deployment or DeploymentConfig, check configuration")
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	result = util.SetResult(resultCode, messages)
	return result
}
