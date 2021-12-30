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
	"github.com/fossul/fossul/src/engine/util"
)

func (s storagePlugin) Backup(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message

	if config.StoragePluginParameters["DeploymentType"] == "DeploymentConfig" {
		result = DeploymentConfigBackupWorkflow(config)
	} else if config.StoragePluginParameters["DeploymentType"] == "Deployment" {
		result = DeploymentBackupWorkflow(config)
	} else if config.StoragePluginParameters["DeploymentType"] == "VirtualMachine" {
		result = VirtualMachineBackupWorkflow(config)
	} else {
		msg := util.SetMessage("ERROR", "CSI storage plugin parameters [DeploymentType] must be DeploymentConfig, Deployment or VirtualMachine")
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	return result
}
