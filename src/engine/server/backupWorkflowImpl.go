package main

import (
	"engine/util"
	"engine/client"
)

func startBackupWorkflowImpl (dataDir string, config util.Config, workflow *util.Workflow) int {
		resultsDir := dataDir + config.ProfileName + "/" + config.ConfigName + "/" + util.IntToString(workflow.Id)
		policy := config.SelectedBackupPolicy
		var isQuiesce bool = false

		commentMsg := "Performing Application Quiesce"
		setComment(resultsDir,commentMsg,workflow)
	
		if config.PreAppQuiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.PreQuiesceCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,step,workflow,result);resultCode != 0 {
				return resultCode
			}
		}
	
		if config.AppQuiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.QuiesceCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,step,workflow,result);resultCode != 0 {
				return resultCode
			}
			isQuiesce = true
		}	
		
		if config.AppPlugin != "" {
			step := stepInit(resultsDir,workflow)
			result := client.Quiesce(config)
			if resultCode := stepErrorHandler(resultsDir,policy,step,workflow,result);resultCode != 0 {
				return resultCode
			}
			isQuiesce = true
		}	
	
		if config.PostAppQuiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.PostQuiesceCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,step,workflow,result);resultCode != 0 {
				unquiesceOnError(resultsDir,policy,isQuiesce,workflow,config)
				return resultCode
			}
		}	
	
		commentMsg = "Performing Backup"
		setComment(resultsDir,commentMsg,workflow)
	
		if config.BackupCreateCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.BackupCreateCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,step,workflow,result);resultCode != 0 {
				unquiesceOnError(resultsDir,policy,isQuiesce,workflow,config)				
				return resultCode
			}
		}	
	
		if config.StoragePlugin != "" {	
			step := stepInit(resultsDir,workflow)
			result := client.Backup(config)
			if resultCode := stepErrorHandler(resultsDir,policy,step,workflow,result);resultCode != 0 {
				unquiesceOnError(resultsDir,policy,isQuiesce,workflow,config)
				return resultCode
			}
		}	

		commentMsg = "Performing Application Unquiesce"
		setComment(resultsDir,commentMsg,workflow)
	
		if config.PreAppUnquiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.PreUnquiesceCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,step,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	
	
		if config.AppUnquiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.UnquiesceCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,step,workflow,result);resultCode != 0 {
				unquiesceOnError(resultsDir,policy,isQuiesce,workflow,config)
				return resultCode
			}
		}	
	
		if config.AppPlugin != "" {
			step := stepInit(resultsDir,workflow)
			result := client.Unquiesce(config)
			if resultCode := stepErrorHandler(resultsDir,policy,step,workflow,result);resultCode != 0 {
				unquiesceOnError(resultsDir,policy,isQuiesce,workflow,config)
				return resultCode
			}
		}	
	
		if config.PostAppUnquiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.PostUnquiesceCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,step,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	

		commentMsg = "Performing Backup Retention"
		setComment(resultsDir,commentMsg,workflow)
	
		if config.BackupDeleteCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.BackupDeleteCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,step,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	
	
	
		if config.StoragePlugin != "" {	
			step := stepInit(resultsDir,workflow)
			result := client.BackupDelete(config)
			if resultCode := stepErrorHandler(resultsDir,policy,step,workflow,result);resultCode != 0 {
				return resultCode
			}
		}	
	
		if config.SendTrapSuccessCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.SendTrapSuccessCmd(config)
			if resultCode := stepErrorHandler(resultsDir,policy,step,workflow,result);resultCode != 0 {
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

	step := util.CreateStep("COMPLETE",workflow)
	util.SetWorkflowStep(workflow,step)
	util.SerializeWorkflow(resultsDir,workflow)
	util.SerializeWorkflowStepResults(resultsDir,step.Id,commentResult)
}

func stepErrorHandler(resultsDir,policy string,step util.Step,workflow *util.Workflow,result util.Result) int {
	if result.Code != 0 {
		workflow.Steps[step.Id].Status="ERROR"
		util.SerializeWorkflowStepResults(resultsDir,step.Id,result)

		util.SetWorkflowStatusError(workflow)
		util.SerializeWorkflow(resultsDir,workflow)

		//remove workflow lock
		delete(runningWorkflowMap,policy)

		return 1
	} else {
		workflow.Steps[step.Id].Status="COMPLETE"
		util.SerializeWorkflowStepResults(resultsDir,step.Id,result)
		util.SerializeWorkflow(resultsDir,workflow)

		return 0
	}
}

func stepInit(resultsDir string,workflow *util.Workflow) util.Step {
	step := util.CreateStep("RUNNING",workflow)
	util.SetWorkflowStep(workflow,step)
	util.SerializeWorkflow(resultsDir,workflow)

	return step
}

func unquiesceOnError(resultsDir,policy string,isQuiesce bool,workflow *util.Workflow,config util.Config) {
	if isQuiesce {
		if config.AppUnquiesceCmd != "" {
			step := stepInit(resultsDir,workflow)
			result := client.UnquiesceCmd(config)
			stepErrorHandler(resultsDir,policy,step,workflow,result)
		}
	
		if config.AppPlugin != "" {
			step := stepInit(resultsDir,workflow)
			result := client.Unquiesce(config)
			stepErrorHandler(resultsDir,policy,step,workflow,result)
		}
	}	
}