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
	"strings"

	"github.com/fossul/fossul/src/client"
	"github.com/fossul/fossul/src/engine/util"
)

func startBackupWorkflowImpl(dataDir string, config util.Config, workflow *util.Workflow) int {
	auth := SetAuth()
	resultsDir := dataDir + "/" + config.ProfileName + "/" + config.ConfigName + "/" + util.IntToString(workflow.Id)
	policy := config.SelectedBackupPolicy
	var isQuiesce bool = false
	var isMount bool = false

	if config.AppPlugin != "" && config.AutoDiscovery == true {
		commentMsg := "Performing Application Discovery"
		setComment(resultsDir, commentMsg, workflow)

		step := stepInit(resultsDir, workflow)

		discoverResult, err := client.Discover(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, discoverResult.Result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, discoverResult.Result, config); resultCode != 0 {
			return resultCode
		}

		// save discovered files in config struct
		if len(config.StoragePluginParameters) == 0 {
			config.StoragePluginParameters = map[string]string{}
		}

		dataFilePaths, logFilePaths := setDiscoverFileList(config, discoverResult)

		dataFilePathsToString := strings.Join(dataFilePaths, ",")
		if len(dataFilePathsToString) != 0 {
			config.StoragePluginParameters["DataFilePaths"] = dataFilePathsToString
		}

		logFilePathsToString := strings.Join(logFilePaths, ",")
		if len(logFilePathsToString) != 0 {
			config.StoragePluginParameters["LogFilePaths"] = logFilePathsToString
		}
	}

	commentMsg := "Performing Application Quiesce"
	setComment(resultsDir, commentMsg, workflow)

	if config.PreAppQuiesceCmd != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.PreQuiesceCmd(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	if config.AppQuiesceCmd != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.QuiesceCmd(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
		isQuiesce = true
	}

	if config.AppPlugin != "" {
		isQuiesce = true
		step := stepInit(resultsDir, workflow)
		result, err := client.Quiesce(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	if config.PostAppQuiesceCmd != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.PostQuiesceCmd(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	commentMsg = "Performing Backup"
	setComment(resultsDir, commentMsg, workflow)

	if config.BackupCreateCmd != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.BackupCreateCmd(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	if config.StoragePlugin != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.Backup(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}

		util.SerializeBackup(resultsDir, &result.Backup)
	}

	commentMsg = "Performing Application Unquiesce"
	setComment(resultsDir, commentMsg, workflow)

	if config.PreAppUnquiesceCmd != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.PreUnquiesceCmd(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	if config.AppUnquiesceCmd != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.UnquiesceCmd(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
		isQuiesce = false
	}

	if config.AppPlugin != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.Unquiesce(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			//unquiesceOnError(resultsDir,policy,isQuiesce,workflow,config)
			return resultCode
		}
		isQuiesce = false
	}

	if config.PostAppUnquiesceCmd != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.PostUnquiesceCmd(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	commentMsg = "Performing Backup Retention"
	setComment(resultsDir, commentMsg, workflow)

	if config.BackupDeleteCmd != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.BackupDeleteCmd(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	if config.StoragePlugin != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.BackupDeleteWorkflow(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	commentMsg = "Performing Mount"
	setComment(resultsDir, commentMsg, workflow)

	if config.StoragePlugin != "" && config.ArchivePlugin != "" {
		isMount = true
		step := stepInit(resultsDir, workflow)
		result, err := client.Mount(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	commentMsg = "Performing Archive"
	setComment(resultsDir, commentMsg, workflow)

	if config.ArchiveCreateCmd != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.ArchiveCreateCmd(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	if config.ArchivePlugin != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.Archive(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	commentMsg = "Performing Archive Retention"
	setComment(resultsDir, commentMsg, workflow)

	if config.ArchiveDeleteCmd != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.ArchiveDeleteCmd(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	if config.ArchivePlugin != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.ArchiveDelete(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	commentMsg = "Performing Unmount"
	setComment(resultsDir, commentMsg, workflow)

	if config.StoragePlugin != "" && config.ArchivePlugin != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.Unmount(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
		isMount = false
	}

	commentMsg = "Job Cleanup"
	setComment(resultsDir, commentMsg, workflow)

	if config.JobRetention != 0 {
		step := stepInit(resultsDir, workflow)
		result := util.DeleteJobs(dataDir, config.ProfileName, config.ConfigName, config.JobRetention)

		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	commentMsg = "Sending Notifications"
	setComment(resultsDir, commentMsg, workflow)

	if config.SendTrapSuccessCmd != "" {
		step := stepInit(resultsDir, workflow)
		result, err := client.SendTrapSuccessCmd(auth, config)
		if err != nil {
			HttpErrorHandlerBackup(err, isQuiesce, isMount, resultsDir, policy, step, workflow, result, config)
			return 1
		}
		if resultCode := StepErrorHandlerBackup(isQuiesce, isMount, resultsDir, policy, step, workflow, result, config); resultCode != 0 {
			return resultCode
		}
	}

	commentMsg = "Backup Completed Successfully"
	setComment(resultsDir, commentMsg, workflow)

	util.SetWorkflowStatusEnd(workflow)
	util.SerializeWorkflow(resultsDir, workflow)

	//remove workflow lock
	runningWorkflowMap.Delete(config.ProfileName + "-" + config.ConfigName)

	return 0
}
