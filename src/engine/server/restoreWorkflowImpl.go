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
	"github.com/fossul/fossul/src/client"
	"github.com/fossul/fossul/src/engine/util"
)

func startRestoreWorkflowImpl(dataDir string, config util.Config, workflow *util.Workflow) int {
	auth := SetAuth()
	resultsDir := dataDir + "/" + config.ProfileName + "/" + config.ConfigName + "/" + util.IntToString(workflow.Id)
	policy := config.SelectedBackupPolicy

	commentMsg := "Retrieving backup contents"
	setComment(resultsDir, commentMsg, workflow)

	step := stepInit(resultsDir, workflow)
	backupContentsDir := dataDir + "/" + config.ProfileName + "/" + config.ConfigName + "/" + util.IntToString(config.SelectedWorkflowId)
	backup := &util.Backup{}
	err := util.ReadGob(backupContentsDir+"/backup", backup)

	if err != nil {
		result := util.SetResultMessage(1, "ERROR", "Couldn't retrieve backup contents for workflow id ["+config.WorkflowId+"] "+err.Error())
		StepErrorHandler(resultsDir, policy, step, workflow, result, config)
		return result.Code
	}

	config.Backup = *backup

	commentMsg = "Performing Application Pre Restore"
	setComment(resultsDir, commentMsg, workflow)

	if config.PreAppRestoreCmd != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.PreAppRestoreCmd(auth, config)
		if err != nil {
			HttpErrorHandler(err, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandler(resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	if config.AppPlugin != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.PreRestore(auth, config)
		if err != nil {
			HttpErrorHandler(err, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandler(resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	commentMsg = "Performing Restore"
	setComment(resultsDir, commentMsg, workflow)

	if config.RestoreCmd != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.RestoreCmd(auth, config)
		if err != nil {
			HttpErrorHandler(err, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandler(resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	if config.StoragePlugin != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.Restore(auth, config)
		if err != nil {
			HttpErrorHandler(err, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandler(resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	commentMsg = "Performing Application Post Restore"
	setComment(resultsDir, commentMsg, workflow)

	if config.PostAppRestoreCmd != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.PostAppRestoreCmd(auth, config)
		if err != nil {
			HttpErrorHandler(err, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandler(resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	if config.AppPlugin != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.PostRestore(auth, config)
		if err != nil {
			HttpErrorHandler(err, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandler(resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	if config.JobRetention != 0 {
		step := stepInit(resultsDir, workflow)
		result := util.DeleteJobs(dataDir, config.ProfileName, config.ConfigName, config.JobRetention)

		if resultCode := StepErrorHandler(resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	commentMsg = "Sending Notifications"
	setComment(resultsDir, commentMsg, workflow)

	if config.SendTrapSuccessCmd != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.SendTrapSuccessCmd(auth, config)
		if err != nil {
			HttpErrorHandler(err, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandler(resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	commentMsg = "Restore Completed Successfully"
	setComment(resultsDir, commentMsg, workflow)

	util.SetWorkflowStatusEnd(workflow)
	util.SerializeWorkflow(resultsDir, workflow)

	//remove workflow lock
	runningWorkflowMap.Delete(config.ProfileName + "-" + config.ConfigName)

	return 0
}
