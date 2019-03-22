package main

import (
	"engine/util"
	"engine/client"
)

func startBackupWorkflowImpl (dataDir string, config util.Config, workflow *util.Workflow) int {
		resultsDir := dataDir + config.ProfileName + "/" + config.ConfigName + "/" + util.IntToString(workflow.Id)
		policy := config.SelectedBackupPolicy

		commentMsg := "Performing Application Quiesce"
		setComment(resultsDir,commentMsg,2,workflow)
	
		if config.PreAppQuiesceCmd != "" {
			stepInit(resultsDir,3,workflow)
			result := client.PreQuiesceCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,3,workflow,result);resultCode != 0 {
				return resultCode
			}
		}
	
		if config.AppQuiesceCmd != "" {
			stepInit(resultsDir,4,workflow)
			result := client.QuiesceCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,4,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	
		
		if config.AppPlugin != "" {
			stepInit(resultsDir,5,workflow)
			result := client.Quiesce(config)
			if resultCode := stepErrorHandler(resultsDir,policy,5,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	
	
		if config.PostAppQuiesceCmd != "" {
			stepInit(resultsDir,6,workflow)
			result := client.PostQuiesceCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,6,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	
	
		commentMsg = "Performing Backup"
		setComment(resultsDir,commentMsg,7,workflow)
	
		if config.BackupCreateCmd != "" {
			stepInit(resultsDir,8,workflow)
			result := client.BackupCreateCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,8,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	
	
		if config.StoragePlugin != "" {	
			stepInit(resultsDir,9,workflow)
			result := client.Backup(config)
			if resultCode := stepErrorHandler(resultsDir,policy,9,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	

		commentMsg = "Performing Application Unquiesce"
		setComment(resultsDir,commentMsg,10,workflow)
	
		if config.PreAppUnquiesceCmd != "" {
			stepInit(resultsDir,11,workflow)
			result := client.PreUnquiesceCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,11,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	
	
		if config.AppUnquiesceCmd != "" {
			stepInit(resultsDir,12,workflow)
			result := client.UnquiesceCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,12,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	
	
		if config.AppPlugin != "" {
			stepInit(resultsDir,13,workflow)
			result := client.Unquiesce(config)
			if resultCode := stepErrorHandler(resultsDir,policy,13,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	
	
		if config.PostAppUnquiesceCmd != "" {
			stepInit(resultsDir,14,workflow)
			result := client.PostUnquiesceCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,14,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	

		commentMsg = "Performing Backup Retention"
		setComment(resultsDir,commentMsg,15,workflow)
	
		if config.BackupDeleteCmd != "" {
			stepInit(resultsDir,16,workflow)
			result := client.BackupDeleteCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,16,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	
	
	
		if config.StoragePlugin != "" {	
			stepInit(resultsDir,17,workflow)
			result := client.BackupDelete(config)
			if resultCode := stepErrorHandler(resultsDir,policy,17,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	
	
		if config.SendTrapSuccessCmd != "" {
			stepInit(resultsDir,18,workflow)
			result := client.SendTrapSuccessCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,18,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	

		commentMsg = "Backup Completed Successfully"
		setComment(resultsDir,commentMsg,19,workflow)

		workflow = util.SetWorkflowStatusEnd(workflow)
		util.SerializeWorkflow(resultsDir,workflow)

		//remove workflow lock
		delete(runningWorkflowMap,config.SelectedBackupPolicy)

		return 0
}

func setComment(resultsDir,msg string,stepId int,workflow *util.Workflow) {
	commentResult := util.SetResultMessage(stepId,"COMMENT",msg)
	//results = append(results, commentResult)
	step := util.SetStep(stepId,"COMPLETE","Step " + util.IntToString(stepId))
	workflow = util.SetWorkflowStep(workflow,step)
	util.SerializeWorkflowStepResults(resultsDir,step.Id,commentResult)
	util.SerializeWorkflow(resultsDir,workflow)
}

func stepErrorHandler(resultsDir,policy string,stepId int,workflow *util.Workflow,result util.Result) int {
	util.SerializeWorkflowStepResults(resultsDir,stepId,result)
	//results = append(results, preQuiesceCmdResult)

	if result.Code != 0 {
		step := util.SetStep(stepId,"ERROR","Step " + util.IntToString(stepId))
		workflow = util.SetWorkflowStep(workflow,step)
		workflow = util.SetWorkflowStatusError(workflow)
		util.SerializeWorkflow(resultsDir,workflow)

		//remove workflow lock
		delete(runningWorkflowMap,policy)

		return 1
	} else {
		step := util.SetStep(stepId,"COMPLETE","Step " + util.IntToString(stepId))
		workflow = util.SetWorkflowStep(workflow,step)
		util.SerializeWorkflow(resultsDir,workflow)
		return 0
	}
}

func stepInit(resultsDir string,stepId int,workflow *util.Workflow) {
	step := util.SetStep(stepId,"RUNNING","Step " + util.IntToString(stepId))
	workflow = util.SetWorkflowStep(workflow,step)
	util.SerializeWorkflow(resultsDir,workflow)
}