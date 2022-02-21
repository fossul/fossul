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

	"github.com/fossul/fossul/src/client"
	"github.com/fossul/fossul/src/engine/util"
)

func startOperatorBackupWorkflowImpl(dataDir string, config util.Config, workflow *util.Workflow) int {
	auth := SetAuth()
	resultsDir := dataDir + "/" + config.ProfileName + "/" + config.ConfigName + "/" + util.IntToString(workflow.Id)
	policy := config.SelectedBackupPolicy
	var isQuiesce bool = false
	var isMount bool = false

	commentMsg := "Performing Backup"
	setComment(resultsDir, commentMsg, workflow)
	step := stepInit(resultsDir, workflow)
	timestampToString := fmt.Sprintf("%d", config.WorkflowTimestamp)
	result, err := client.CreateCustomBackupResource(auth, config.ProfileName, config.ConfigName, policy, config.WorkflowId, timestampToString)
	workflow = setLastMessage(result, workflow)
	if err != nil {
		HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
		return 1
	}

	if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
		return resultCode
	}

	commentMsg = "Performing Backup Retention"
	setComment(resultsDir, commentMsg, workflow)
	step = stepInit(resultsDir, workflow)
	result, err = client.BackupCustomResourceRetention(auth, config.ProfileName, config.ConfigName, policy)
	workflow = setLastMessage(result, workflow)

	if err != nil {
		HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
		return 1
	}
	if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
		return resultCode
	}

	commentMsg = "Backup Workflow Completed Successfully"
	setComment(resultsDir, commentMsg, workflow)

	util.SetWorkflowStatusEnd(workflow)
	util.SerializeWorkflow(resultsDir, workflow)

	return 0
}
