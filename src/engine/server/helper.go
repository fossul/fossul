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
	"fossul/src/engine/client"
	"fossul/src/engine/util"
)

func setComment(resultsDir, msg string, workflow *util.Workflow) {
	commentResult := util.SetResultMessage(0, "COMMENT", msg)

	step := util.CreateCommentStep(workflow)
	util.SetWorkflowStep(workflow, step)
	util.SerializeWorkflow(resultsDir, workflow)
	util.SerializeWorkflowStepResults(resultsDir, step.Id, commentResult)
}

func StepErrorHandlerBackup(isQuiesce, isMount bool, resultsDir, policy string, step util.Step, workflow *util.Workflow, result util.Result, config util.Config) int {
	auth := SetAuth()

	if result.Code != 0 {
		util.SetStepError(workflow, step)
		util.SerializeWorkflowStepResults(resultsDir, step.Id, result)

		if isQuiesce {
			commentMsg := "Performing Application Unquiesce"
			setComment(resultsDir, commentMsg, workflow)

			if config.AppUnquiesceCmd != "" {
				step := stepInit(resultsDir, workflow)
				result, _ := client.UnquiesceCmd(auth, config)

				util.SetStepError(workflow, step)
				util.SerializeWorkflowStepResults(resultsDir, step.Id, result)
			}

			if config.AppPlugin != "" {
				step := stepInit(resultsDir, workflow)
				result, _ := client.Unquiesce(auth, config)

				util.SetStepError(workflow, step)
				util.SerializeWorkflowStepResults(resultsDir, step.Id, result)
			}
		}

		if isMount {
			commentMsg := "Performing Unmount"
			setComment(resultsDir, commentMsg, workflow)

			if config.StoragePlugin != "" && config.ArchivePlugin != "" {
				step := stepInit(resultsDir, workflow)
				result, _ := client.Unmount(auth, config)

				util.SetStepError(workflow, step)
				util.SerializeWorkflowStepResults(resultsDir, step.Id, result)
			}
		}

		sendErrorNotification(resultsDir, policy, step, workflow, result, config)
		util.SetWorkflowStatusError(workflow)
		util.SerializeWorkflow(resultsDir, workflow)

		//remove workflow lock
		delete(runningWorkflowMap, config.ProfileName+"-"+config.ConfigName)

		return 1
	} else {
		util.SetStepComplete(workflow, step)
		util.SerializeWorkflowStepResults(resultsDir, step.Id, result)
		util.SerializeWorkflow(resultsDir, workflow)

		return 0
	}
}

func StepErrorHandler(resultsDir, policy string, step util.Step, workflow *util.Workflow, result util.Result, config util.Config) int {

	if result.Code != 0 {
		util.SetStepError(workflow, step)
		util.SerializeWorkflowStepResults(resultsDir, step.Id, result)

		sendErrorNotification(resultsDir, policy, step, workflow, result, config)
		util.SetWorkflowStatusError(workflow)
		util.SerializeWorkflow(resultsDir, workflow)

		//remove workflow lock
		delete(runningWorkflowMap, config.ProfileName+"-"+config.ConfigName)

		return 1
	} else {
		util.SetStepComplete(workflow, step)
		util.SerializeWorkflowStepResults(resultsDir, step.Id, result)
		util.SerializeWorkflow(resultsDir, workflow)

		return 0
	}
}

func HttpErrorHandlerBackup(err error, isQuiesce, isMount bool, resultsDir, policy string, step util.Step, workflow *util.Workflow, result util.Result, config util.Config) {
	auth := SetAuth()

	msg := util.SetMessage("ERROR", err.Error())
	result.Messages = util.PrependMessage(msg, result.Messages)
	result.Code = 1

	util.SetStepError(workflow, step)
	util.SerializeWorkflowStepResults(resultsDir, step.Id, result)

	if isQuiesce {
		commentMsg := "Performing Application Unquiesce"
		setComment(resultsDir, commentMsg, workflow)

		if config.AppUnquiesceCmd != "" {
			step := stepInit(resultsDir, workflow)
			result, _ := client.UnquiesceCmd(auth, config)
			util.SetStepError(workflow, step)
			util.SerializeWorkflowStepResults(resultsDir, step.Id, result)
		}

		if config.AppPlugin != "" {
			step := stepInit(resultsDir, workflow)
			result, _ := client.Unquiesce(auth, config)
			util.SetStepError(workflow, step)
			util.SerializeWorkflowStepResults(resultsDir, step.Id, result)
		}
	}

	if isMount {
		commentMsg := "Performing Unmount"
		setComment(resultsDir, commentMsg, workflow)

		if config.StoragePlugin != "" && config.ArchivePlugin != "" {
			step := stepInit(resultsDir, workflow)
			result, _ := client.Unmount(auth, config)

			util.SetStepError(workflow, step)
			util.SerializeWorkflowStepResults(resultsDir, step.Id, result)
		}
	}

	sendErrorNotification(resultsDir, policy, step, workflow, result, config)
	util.SetWorkflowStatusError(workflow)
	util.SerializeWorkflow(resultsDir, workflow)

	//remove workflow lock
	delete(runningWorkflowMap, config.ProfileName+"-"+config.ConfigName)
}

func HttpErrorHandler(err error, resultsDir, policy string, step util.Step, workflow *util.Workflow, result util.Result, config util.Config) {

	msg := util.SetMessage("ERROR", err.Error())
	result.Messages = util.PrependMessage(msg, result.Messages)
	result.Code = 1

	util.SetStepError(workflow, step)
	util.SerializeWorkflowStepResults(resultsDir, step.Id, result)

	sendErrorNotification(resultsDir, policy, step, workflow, result, config)
	util.SetWorkflowStatusError(workflow)
	util.SerializeWorkflow(resultsDir, workflow)

	//remove workflow lock
	delete(runningWorkflowMap, config.ProfileName+"-"+config.ConfigName)
}

func stepInit(resultsDir string, workflow *util.Workflow) util.Step {
	step := util.CreateStep(workflow)
	util.SetWorkflowStep(workflow, step)
	util.SerializeWorkflow(resultsDir, workflow)

	return step
}

func sendErrorNotification(resultsDir, policy string, step util.Step, workflow *util.Workflow, result util.Result, config util.Config) {
	auth := SetAuth()

	if config.SendTrapErrorCmd != "" {
		commentMsg := "Sending Error Notifications"
		setComment(resultsDir, commentMsg, workflow)

		step := stepInit(resultsDir, workflow)
		result, _ := client.SendTrapErrorCmd(auth, config)

		if result.Code != 0 {
			util.SetStepError(workflow, step)
		} else {
			util.SetStepComplete(workflow, step)
		}

		util.SerializeWorkflowStepResults(resultsDir, step.Id, result)
		util.SerializeWorkflow(resultsDir, workflow)
	}
}

func setDiscoverFileList(config util.Config, discoverResult util.DiscoverResult) (dataFilePaths, logFilePaths []string) {
	for _, discover := range discoverResult.DiscoverList {
		for _, dataFilePath := range discover.DataFilePaths {
			dataFilePaths = append(dataFilePaths, dataFilePath)
		}

		for _, logFilePath := range discover.LogFilePaths {
			logFilePaths = append(logFilePaths, logFilePath)
		}
	}

	return dataFilePaths, logFilePaths
}

func SetAuth() client.Auth {
	var auth client.Auth
	auth.ServerHostname = serverHostname
	auth.ServerPort = serverPort
	auth.AppHostname = appHostname
	auth.AppPort = appPort
	auth.StorageHostname = storageHostname
	auth.StoragePort = storagePort
	auth.Username = myUser
	auth.Password = myPass

	return auth
}
