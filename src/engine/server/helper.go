package main

import (
	"fossil/src/engine/util"
	"fossil/src/engine/client"
)

func setComment(resultsDir,msg string,workflow *util.Workflow)  {
	commentResult := util.SetResultMessage(0,"COMMENT",msg)

	step := util.CreateCommentStep(workflow)
	util.SetWorkflowStep(workflow,step)
	util.SerializeWorkflow(resultsDir,workflow)
	util.SerializeWorkflowStepResults(resultsDir,step.Id,commentResult)
}

func StepErrorHandlerBackup(isQuiesce bool,resultsDir,policy string,step util.Step,workflow *util.Workflow,result util.Result,config util.Config) int {
	auth := SetAuth()

	if result.Code != 0 {
		util.SetStepError(workflow,step)
		util.SerializeWorkflowStepResults(resultsDir,step.Id,result)

		if isQuiesce {
			commentMsg := "Performing Application Unquiesce"
			setComment(resultsDir,commentMsg,workflow)
	
			if config.AppUnquiesceCmd != "" {
				step := stepInit(resultsDir,workflow)
				result,_ := client.UnquiesceCmd(auth,config)

				util.SetStepError(workflow,step)
				util.SerializeWorkflowStepResults(resultsDir,step.Id,result)
			}
		
			if config.AppPlugin != "" {
				step := stepInit(resultsDir,workflow)
				result,_ := client.Unquiesce(auth,config)

				util.SetStepError(workflow,step)
				util.SerializeWorkflowStepResults(resultsDir,step.Id,result)
			}
		}	

		sendErrorNotification(resultsDir,policy,step,workflow,result,config)
		util.SetWorkflowStatusError(workflow)
		util.SerializeWorkflow(resultsDir,workflow)

		//remove workflow lock
		delete(runningWorkflowMap,config.ProfileName + "-" + config.ConfigName)

		return 1
	} else {
		util.SetStepComplete(workflow,step)
		util.SerializeWorkflowStepResults(resultsDir,step.Id,result)
		util.SerializeWorkflow(resultsDir,workflow)

		return 0
	}
}

func StepErrorHandler(resultsDir,policy string,step util.Step,workflow *util.Workflow,result util.Result,config util.Config) int {
	
	if result.Code != 0 {
		util.SetStepError(workflow,step)
		util.SerializeWorkflowStepResults(resultsDir,step.Id,result)

		sendErrorNotification(resultsDir,policy,step,workflow,result,config)
		util.SetWorkflowStatusError(workflow)
		util.SerializeWorkflow(resultsDir,workflow)

		//remove workflow lock
		delete(runningWorkflowMap,config.ProfileName + "-" + config.ConfigName)

		return 1
	} else {
		util.SetStepComplete(workflow,step)
		util.SerializeWorkflowStepResults(resultsDir,step.Id,result)
		util.SerializeWorkflow(resultsDir,workflow)

		return 0
	}
}

func HttpErrorHandlerBackup(err error,isQuiesce bool,resultsDir,policy string,step util.Step,workflow *util.Workflow,result util.Result,config util.Config) {
	auth := SetAuth()

	msg := util.SetMessage("ERROR",err.Error())
	result.Messages = util.PrependMessage(msg,result.Messages)
	result.Code = 1
	
	util.SetStepError(workflow,step)
	util.SerializeWorkflowStepResults(resultsDir,step.Id,result)

	if isQuiesce {
		commentMsg := "Performing Application Unquiesce"
		setComment(resultsDir,commentMsg,workflow)
	
		if config.AppUnquiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result,_ := client.UnquiesceCmd(auth,config)
			util.SetStepError(workflow,step)
			util.SerializeWorkflowStepResults(resultsDir,step.Id,result)
		}
		
		if config.AppPlugin != "" {
			step := stepInit(resultsDir,workflow)
			result,_ := client.Unquiesce(auth,config)
			util.SetStepError(workflow,step)
			util.SerializeWorkflowStepResults(resultsDir,step.Id,result)
		}
	}	

	sendErrorNotification(resultsDir,policy,step,workflow,result,config)
	util.SetWorkflowStatusError(workflow)
	util.SerializeWorkflow(resultsDir,workflow)

	//remove workflow lock
	delete(runningWorkflowMap,config.ProfileName + "-" + config.ConfigName)
}

func HttpErrorHandler(err error,resultsDir,policy string,step util.Step,workflow *util.Workflow,result util.Result,config util.Config) {

	msg := util.SetMessage("ERROR",err.Error())
	result.Messages = util.PrependMessage(msg,result.Messages)
	result.Code = 1
	
	util.SetStepError(workflow,step)
	util.SerializeWorkflowStepResults(resultsDir,step.Id,result)

	sendErrorNotification(resultsDir,policy,step,workflow,result,config)
	util.SetWorkflowStatusError(workflow)
	util.SerializeWorkflow(resultsDir,workflow)

	//remove workflow lock
	delete(runningWorkflowMap,config.ProfileName + "-" + config.ConfigName)
}

func stepInit(resultsDir string,workflow *util.Workflow) util.Step {
	step := util.CreateStep(workflow)
	util.SetWorkflowStep(workflow,step)
	util.SerializeWorkflow(resultsDir,workflow)

	return step
}

func sendErrorNotification(resultsDir,policy string,step util.Step,workflow *util.Workflow,result util.Result,config util.Config) {
	auth := SetAuth()

	if config.SendTrapErrorCmd != "" {
		commentMsg := "Sending Error Notifications"
		setComment(resultsDir,commentMsg,workflow)	

		step := stepInit(resultsDir,workflow)
		result,_ := client.SendTrapErrorCmd(auth,config)

		if result.Code != 0 {
			util.SetStepError(workflow,step)
		} else {
			util.SetStepComplete(workflow,step)		
		}

		util.SerializeWorkflowStepResults(resultsDir,step.Id,result)
		util.SerializeWorkflow(resultsDir,workflow)
	}	
}

func setDiscoverFileList(config util.Config, discoverResult util.DiscoverResult) (dataFilePaths,logFilePaths []string) {
	for _,discover := range discoverResult.DiscoverList {
		for _,dataFilePath := range discover.DataFilePaths {
			dataFilePaths = append(dataFilePaths,dataFilePath)
		}

		for _,logFilePath := range discover.LogFilePaths {
			logFilePaths = append(logFilePaths,logFilePath)
		}
	}

	return dataFilePaths,logFilePaths
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