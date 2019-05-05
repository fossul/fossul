package main

import (
	"fossil/src/engine/util"
	"fossil/src/engine/client"
	"strings"
)

func startBackupWorkflowImpl (dataDir string, config util.Config, workflow *util.Workflow) int {
	auth := setAuth()
	resultsDir := dataDir + config.ProfileName + "/" + config.ConfigName + "/" + util.IntToString(workflow.Id)
	policy := config.SelectedBackupPolicy
	var isQuiesce bool = false
	
	if config.AppPlugin != "" && config.AutoDiscovery == true {
		commentMsg := "Performing Application Discovery"
		setComment(resultsDir,commentMsg,workflow)

		step := stepInit(resultsDir,workflow)
			
		discoverResult,err := client.Discover(auth,config)
		if err != nil {
			httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,discoverResult.Result,config)
			return 1
		}
		if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,discoverResult.Result,config);resultCode != 0 {
			return resultCode
		}

		// save discovered files in config struct
		dataFilePaths,logFilePaths := setDiscoverFileList(config,discoverResult)
			
		dataFilePathsToString := strings.Join(dataFilePaths,",")
		config.StoragePluginParameters["DataFilePaths"] = dataFilePathsToString

		logFilePathsToString := strings.Join(logFilePaths,",")
		config.StoragePluginParameters["LogFilePaths"] = logFilePathsToString
	}	

		commentMsg := "Performing Application Quiesce"
		setComment(resultsDir,commentMsg,workflow)
	
		if config.PreAppQuiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result,err := client.PreQuiesceCmd(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}
	
		if config.AppQuiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result,err := client.QuiesceCmd(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
			isQuiesce = true
		}	
		
		if config.AppPlugin != "" {
			isQuiesce = true
			step := stepInit(resultsDir,workflow)
			result,err := client.Quiesce(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	
	
		if config.PostAppQuiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result,err := client.PostQuiesceCmd(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	
	
		commentMsg = "Performing Backup"
		setComment(resultsDir,commentMsg,workflow)
	
		if config.BackupCreateCmd != "" {
			step := stepInit(resultsDir,workflow)
			result,err := client.BackupCreateCmd(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {			
				return resultCode
			}
		}	
	
		if config.StoragePlugin != "" {	
			step := stepInit(resultsDir,workflow)
			result,err := client.Backup(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	

		commentMsg = "Performing Application Unquiesce"
		setComment(resultsDir,commentMsg,workflow)
	
		if config.PreAppUnquiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result,err := client.PreUnquiesceCmd(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	
	
		if config.AppUnquiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result,err := client.UnquiesceCmd(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
			isQuiesce = false
		}
	
		if config.AppPlugin != "" {
			step := stepInit(resultsDir,workflow)
			result,err := client.Unquiesce(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				//unquiesceOnError(resultsDir,policy,isQuiesce,workflow,config)
				return resultCode
			}
			isQuiesce = false
		}	

		if config.PostAppUnquiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result,err := client.PostUnquiesceCmd(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	

		commentMsg = "Performing Backup Retention"
		setComment(resultsDir,commentMsg,workflow)
	
		if config.BackupDeleteCmd != "" {
			step := stepInit(resultsDir,workflow)
			result,err := client.BackupDeleteCmd(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	
	
	
		if config.StoragePlugin != "" {	
			step := stepInit(resultsDir,workflow)
			result,err := client.BackupDelete(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}

		commentMsg = "Performing Archive Retention"
		setComment(resultsDir,commentMsg,workflow)

		if config.ArchiveCreateCmd != "" {
			step := stepInit(resultsDir,workflow)
			result,err := client.ArchiveCreateCmd(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	
		
		if config.ArchivePlugin != "" {	
			step := stepInit(resultsDir,workflow)
			result,err := client.Archive(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}

		if config.ArchiveDeleteCmd != "" {
			step := stepInit(resultsDir,workflow)
			result,err := client.ArchiveDeleteCmd(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	
		
		if config.ArchivePlugin != "" {	
			step := stepInit(resultsDir,workflow)
			result,err := client.ArchiveDelete(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}

		commentMsg = "Sending Notifications"
		setComment(resultsDir,commentMsg,workflow)	
	
		if config.SendTrapSuccessCmd != "" {
			step := stepInit(resultsDir,workflow)
			result,err := client.SendTrapSuccessCmd(auth,config)
			if err != nil {
				httpErrorHandler(err,isQuiesce,resultsDir,policy,step,workflow,result,config)
				return 1
			}
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	

		commentMsg = "Backup Completed Successfully"
		setComment(resultsDir,commentMsg,workflow)

		util.SetWorkflowStatusEnd(workflow)
		util.SerializeWorkflow(resultsDir,workflow)

		//remove workflow lock
		delete(runningWorkflowMap,config.SelectedBackupPolicy)

		return 0
}

func setComment(resultsDir,msg string,workflow *util.Workflow)  {
	commentResult := util.SetResultMessage(0,"COMMENT",msg)

	step := util.CreateCommentStep(workflow)
	util.SetWorkflowStep(workflow,step)
	util.SerializeWorkflow(resultsDir,workflow)
	util.SerializeWorkflowStepResults(resultsDir,step.Id,commentResult)
}

func stepErrorHandler(isQuiesce bool,resultsDir,policy string,step util.Step,workflow *util.Workflow,result util.Result,config util.Config) int {
	auth := setAuth()

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
		delete(runningWorkflowMap,policy)

		return 1
	} else {
		util.SetStepComplete(workflow,step)
		util.SerializeWorkflowStepResults(resultsDir,step.Id,result)
		util.SerializeWorkflow(resultsDir,workflow)

		return 0
	}
}

func httpErrorHandler(err error,isQuiesce bool,resultsDir,policy string,step util.Step,workflow *util.Workflow,result util.Result,config util.Config) {
	auth := setAuth()

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
	delete(runningWorkflowMap,policy)
}

func stepInit(resultsDir string,workflow *util.Workflow) util.Step {
	step := util.CreateStep(workflow)
	util.SetWorkflowStep(workflow,step)
	util.SerializeWorkflow(resultsDir,workflow)

	return step
}

func sendErrorNotification(resultsDir,policy string,step util.Step,workflow *util.Workflow,result util.Result,config util.Config) {
	auth := setAuth()

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

func setAuth() client.Auth {
	var auth client.Auth
	auth.Username = myUser
	auth.Password = myPass

	return auth
}
