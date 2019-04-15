package main

import (
	"engine/util"
	"engine/client"
	"strings"
)

func startBackupWorkflowImpl (dataDir string, config util.Config, workflow *util.Workflow) int {
		resultsDir := dataDir + config.ProfileName + "/" + config.ConfigName + "/" + util.IntToString(workflow.Id)
		policy := config.SelectedBackupPolicy
		var isQuiesce bool = false

		if config.AppPlugin != "" && config.AutoDiscovery == true {
			commentMsg := "Performing Application Discovery"
			setComment(resultsDir,commentMsg,workflow)

			step := stepInit(resultsDir,workflow)
			discoverResult := client.Discover(config)
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
			result := client.PreQuiesceCmd(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}
	
		if config.AppQuiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.QuiesceCmd(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
			isQuiesce = true
		}	
		
		if config.AppPlugin != "" {
			isQuiesce = true
			step := stepInit(resultsDir,workflow)
			result := client.Quiesce(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	
	
		if config.PostAppQuiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.PostQuiesceCmd(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	
	
		commentMsg = "Performing Backup"
		setComment(resultsDir,commentMsg,workflow)
	
		if config.BackupCreateCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.BackupCreateCmd(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {			
				return resultCode
			}
		}	
	
		if config.StoragePlugin != "" {	
			step := stepInit(resultsDir,workflow)
			result := client.Backup(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	

		commentMsg = "Performing Application Unquiesce"
		setComment(resultsDir,commentMsg,workflow)
	
		if config.PreAppUnquiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.PreUnquiesceCmd(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	
	
		if config.AppUnquiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.UnquiesceCmd(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
			isQuiesce = false
		}	
	
		if config.AppPlugin != "" {
			step := stepInit(resultsDir,workflow)
			result := client.Unquiesce(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				//unquiesceOnError(resultsDir,policy,isQuiesce,workflow,config)
				return resultCode
			}
			isQuiesce = false
		}	

		if config.PostAppUnquiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.PostUnquiesceCmd(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	

		commentMsg = "Performing Backup Retention"
		setComment(resultsDir,commentMsg,workflow)
	
		if config.BackupDeleteCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.BackupDeleteCmd(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	
	
	
		if config.StoragePlugin != "" {	
			step := stepInit(resultsDir,workflow)
			result := client.BackupDelete(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}

		commentMsg = "Performing Archive Retention"
		setComment(resultsDir,commentMsg,workflow)

		if config.ArchiveCreateCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.ArchiveCreateCmd(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	
		
		if config.ArchivePlugin != "" {	
			step := stepInit(resultsDir,workflow)
			result := client.Archive(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}

		if config.ArchiveDeleteCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.ArchiveDeleteCmd(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}	
		
		if config.ArchivePlugin != "" {	
			step := stepInit(resultsDir,workflow)
			result := client.ArchiveDelete(config)
			if resultCode := stepErrorHandler(isQuiesce,resultsDir,policy,step,workflow,result,config);resultCode != 0 {
				return resultCode
			}
		}

		commentMsg = "Sending Notifications"
		setComment(resultsDir,commentMsg,workflow)	
	
		if config.SendTrapSuccessCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.SendTrapSuccessCmd(config)
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
	if result.Code != 0 {
		util.SetStepError(workflow,step)
		util.SerializeWorkflowStepResults(resultsDir,step.Id,result)

		if isQuiesce {
			commentMsg := "Performing Application Unquiesce"
			setComment(resultsDir,commentMsg,workflow)
	
			if config.AppUnquiesceCmd != "" {
				step := stepInit(resultsDir,workflow)
				result := client.UnquiesceCmd(config)
				util.SetStepError(workflow,step)
				util.SerializeWorkflowStepResults(resultsDir,step.Id,result)
			}
		
			if config.AppPlugin != "" {
				step := stepInit(resultsDir,workflow)
				result := client.Unquiesce(config)
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

func stepInit(resultsDir string,workflow *util.Workflow) util.Step {
	step := util.CreateStep(workflow)
	util.SetWorkflowStep(workflow,step)
	util.SerializeWorkflow(resultsDir,workflow)

	return step
}

func sendErrorNotification(resultsDir,policy string,step util.Step,workflow *util.Workflow,result util.Result,config util.Config) {
	if config.SendTrapErrorCmd != "" {
		commentMsg := "Sending Error Notifications"
		setComment(resultsDir,commentMsg,workflow)	

		step := stepInit(resultsDir,workflow)
		result := client.SendTrapErrorCmd(config)

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
