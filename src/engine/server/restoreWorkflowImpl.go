package main

import (
	"fossul/src/engine/util"
	"fossul/src/engine/client"
)

func startRestoreWorkflowImpl (dataDir string, config util.Config, workflow *util.Workflow) int {
	auth := SetAuth()
	resultsDir := dataDir + "/" + config.ProfileName + "/" + config.ConfigName + "/" + util.IntToString(workflow.Id)
	policy := config.SelectedBackupPolicy

	commentMsg := "Performing Application Pre Restore"
	setComment(resultsDir,commentMsg,workflow)

	if config.PreAppRestoreCmd != "" {
		step := stepInit(resultsDir,workflow)
		result,err := client.PreAppRestoreCmd(auth,config)
		if err != nil {
			HttpErrorHandler(err,resultsDir,policy,step,workflow,result,config)
			return 1
		}
		if resultCode := StepErrorHandler(resultsDir,policy,step,workflow,result,config);resultCode != 0 {
			return resultCode
		}
	}

	if config.AppPlugin != "" {
		step := stepInit(resultsDir,workflow)
		result,err := client.PreRestore(auth,config)
		if err != nil {
			HttpErrorHandler(err,resultsDir,policy,step,workflow,result,config)
			return 1
		}
		if resultCode := StepErrorHandler(resultsDir,policy,step,workflow,result,config);resultCode != 0 {
			return resultCode
		}
	}	

	commentMsg = "Performing Restore"
	setComment(resultsDir,commentMsg,workflow)
	
	if config.RestoreCmd != "" {
		step := stepInit(resultsDir,workflow)
		result,err := client.RestoreCmd(auth,config)
		if err != nil {
			HttpErrorHandler(err,resultsDir,policy,step,workflow,result,config)
			return 1
		}
		if resultCode := StepErrorHandler(resultsDir,policy,step,workflow,result,config);resultCode != 0 {
			return resultCode
		}
	}

	if config.StoragePlugin != "" {	
		step := stepInit(resultsDir,workflow)
		result,err := client.Restore(auth,config)
		if err != nil {
			HttpErrorHandler(err,resultsDir,policy,step,workflow,result,config)
			return 1
		}
		if resultCode := StepErrorHandler(resultsDir,policy,step,workflow,result,config);resultCode != 0 {
			return resultCode
		}
	}	

	commentMsg = "Performing Application Post Restore"
	setComment(resultsDir,commentMsg,workflow)
	
	if config.PostAppRestoreCmd != "" {
		step := stepInit(resultsDir,workflow)
		result,err := client.PostAppRestoreCmd(auth,config)
		if err != nil {
			HttpErrorHandler(err,resultsDir,policy,step,workflow,result,config)
			return 1
		}
		if resultCode := StepErrorHandler(resultsDir,policy,step,workflow,result,config);resultCode != 0 {
			return resultCode
		}
	}

	if config.AppPlugin != "" {
		step := stepInit(resultsDir,workflow)
		result,err := client.PostRestore(auth,config)
		if err != nil {
			HttpErrorHandler(err,resultsDir,policy,step,workflow,result,config)
			return 1
		}
		if resultCode := StepErrorHandler(resultsDir,policy,step,workflow,result,config);resultCode != 0 {
			return resultCode
		}
	}	
	
	if config.JobRetention != 0 {
		step := stepInit(resultsDir,workflow)
		result := util.DeleteJobs(dataDir,config.ProfileName,config.ConfigName,config.JobRetention)

		if resultCode := StepErrorHandler(resultsDir,policy,step,workflow,result,config);resultCode != 0 {
			return resultCode
		}
	}

	commentMsg = "Sending Notifications"
	setComment(resultsDir,commentMsg,workflow)	
	
	if config.SendTrapSuccessCmd != "" {
		step := stepInit(resultsDir,workflow)
		result,err := client.SendTrapSuccessCmd(auth,config)
		if err != nil {
			HttpErrorHandler(err,resultsDir,policy,step,workflow,result,config)
			return 1
		}
		if resultCode := StepErrorHandler(resultsDir,policy,step,workflow,result,config);resultCode != 0 {
			return resultCode
		}
	}	

	commentMsg = "Restore Completed Successfully"
	setComment(resultsDir,commentMsg,workflow)

	util.SetWorkflowStatusEnd(workflow)
	util.SerializeWorkflow(resultsDir,workflow)

	//remove workflow lock
	delete(runningWorkflowMap,config.ProfileName + "-" + config.ConfigName)

	return 0
}
