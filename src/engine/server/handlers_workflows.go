package main

import (
	"github.com/gorilla/mux"
	"encoding/json"
	"fossil/src/engine/util"
	"net/http"
	"log"
	"time"
)

func StartBackupWorkflowLocalConfig(w http.ResponseWriter, r *http.Request) {
	var workflowResult util.WorkflowResult
	workflow := &util.Workflow{}
	workflow.Id =  util.GetWorkflowId()
	workflow.Type = "backup"
	workflow.Status = "RUNNING"

	var timestamp string = time.Now().Format(time.RFC3339)
	workflow.Timestamp = timestamp

	workflowResult.Id = workflow.Id

	config,_ := util.GetConfig(w,r)
	config.WorkflowId = util.IntToString(workflow.Id)

	_,ok := runningWorkflowMap[config.ProfileName + "-" + config.ConfigName]
	if ok {	
		result := util.SetResultMessage(1,"ERROR","Workflow id [" + util.IntToString(workflow.Id) + "] failed to start. Another workflow is running under profile [" + config.ProfileName + "] config [" + config.ConfigName + "]")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)		
	} else {
		runningWorkflowMap[config.ProfileName + "-" + config.ConfigName] = config.SelectedBackupPolicy

		go func() {
			startBackupWorkflowImpl(dataDir,config,workflow)
		}()
	
		result := util.SetResultMessage(0,"INFO","Workflow id [" + util.IntToString(workflow.Id) + "] started successfully")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)
	}
}


func StartBackupWorkflow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var policyName string = params["policy"]

	var workflowResult util.WorkflowResult
	workflow := &util.Workflow{}
	workflow.Id =  util.GetWorkflowId()
	workflow.Type = "backup"
	workflow.Status = "RUNNING"

	var timestamp string = time.Now().Format(time.RFC3339)
	workflow.Timestamp = timestamp

	workflowResult.Id = workflow.Id

	config,err := GetConsolidatedConfig(profileName,configName,policyName)
	if err != nil {
		result := util.SetResultMessage(1,"ERROR","Workflow id [" + util.IntToString(workflow.Id) + "] failed to start. Couldn't read config using profile [" + profileName + "] config [" + configName + "] " + err.Error())
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)	
	}

	config.WorkflowId = util.IntToString(workflow.Id)

	_,ok := runningWorkflowMap[config.ProfileName + "-" + config.ConfigName]
	if ok {	
		result := util.SetResultMessage(1,"ERROR","Workflow id [" + util.IntToString(workflow.Id) + "] failed to start. Another workflow is running under profile [" + config.ProfileName + "] config [" + config.ConfigName + "] " + err.Error())
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)		
	} else {
		runningWorkflowMap[config.ProfileName + "-" + config.ConfigName] = config.SelectedBackupPolicy

		go func() {
			startBackupWorkflowImpl(dataDir,config,workflow)
		}()
	
		result := util.SetResultMessage(0,"INFO","Workflow id [" + util.IntToString(workflow.Id) + "] started successfully")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)
	}
}

func StartRestoreWorkflowLocalConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var selectedWorkflowId string = params["workflowId"]

	var workflowResult util.WorkflowResult
	workflow := &util.Workflow{}
	workflow.Id =  util.GetWorkflowId()
	workflow.Type = "restore"
	workflow.Status = "RUNNING"

	var timestamp string = time.Now().Format(time.RFC3339)
	workflow.Timestamp = timestamp

	workflowResult.Id = workflow.Id

	config,_ := util.GetConfig(w,r)
	config.WorkflowId = util.IntToString(workflow.Id)
	config.SelectedWorkflowId = util.StringToInt(selectedWorkflowId)

	_,ok := runningWorkflowMap[config.ProfileName + "-" + config.ConfigName]
	if ok {	
		result := util.SetResultMessage(1,"ERROR","Workflow id [" + util.IntToString(workflow.Id) + "] failed to start. Another workflow is running under profile [" + config.ProfileName + "] config [" + config.ConfigName + "]")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)		
	} else {
		runningWorkflowMap[config.ProfileName + "-" + config.ConfigName] = config.SelectedBackupPolicy

		go func() {
			startRestoreWorkflowImpl(dataDir,config,workflow)
		}()
	
		result := util.SetResultMessage(0,"INFO","Workflow id [" + util.IntToString(workflow.Id) + "] started successfully")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)
	}
}

func StartRestoreWorkflow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var policyName string = params["policy"]
	var selectedWorkflowId string = params["workflowId"]

	var workflowResult util.WorkflowResult
	workflow := &util.Workflow{}
	workflow.Id =  util.GetWorkflowId()
	workflow.Type = "restore"
	workflow.Status = "RUNNING"

	var timestamp string = time.Now().Format(time.RFC3339)
	workflow.Timestamp = timestamp

	workflowResult.Id = workflow.Id

	config,err := GetConsolidatedConfig(profileName,configName,policyName)
	if err != nil {
		result := util.SetResultMessage(1,"ERROR","Workflow id [" + util.IntToString(workflow.Id) + "] failed to start. Couldn't read config using profile [" + profileName + "] config [" + configName + "]")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)	
	}

	config.WorkflowId = util.IntToString(workflow.Id)
	config.SelectedWorkflowId = util.StringToInt(selectedWorkflowId)

	_,ok := runningWorkflowMap[config.ProfileName + "-" + config.ConfigName]
	if ok {	
		result := util.SetResultMessage(1,"ERROR","Workflow id [" + util.IntToString(workflow.Id) + "] failed to start. Another workflow is running under profile [" + config.ProfileName + "] config [" + config.ConfigName + "]")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)		
	} else {
		runningWorkflowMap[config.ProfileName + "-" + config.ConfigName] = config.SelectedBackupPolicy

		go func() {
			startRestoreWorkflowImpl(dataDir,config,workflow)
		}()
	
		result := util.SetResultMessage(0,"INFO","Workflow id [" + util.IntToString(workflow.Id) + "] started successfully")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)
	}
}

func GetWorkflowStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var id string = params["id"]

	resultsDir := dataDir + "/" + profileName + "/" + configName + "/" + id

	workflow := &util.Workflow{}

	workflowFile := resultsDir + "/workflow"
	err := util.ReadGob(workflowFile,&workflow)
	if err != nil {
		log.Println(err.Error())
	}

	_ = json.NewDecoder(r.Body).Decode(&workflow)
	json.NewEncoder(w).Encode(workflow)	
}

func GetWorkflowStepResults(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var workflowId string = params["workflowId"]
	var stepId string = params["stepId"]

	var results []util.Result	
	resultsFile := dataDir + "/" + profileName + "/" + configName + "/" + workflowId + "/" + stepId

	var result util.Result
	err := util.ReadGob(resultsFile,&result)
	results = append(results, result)
	if err != nil {
		log.Println(err.Error())

		var results []util.Result
		errorResult := util.SetResultMessage(1,"ERROR",err.Error())
		results = append(results, errorResult)
	
		_ = json.NewDecoder(r.Body).Decode(&results)
		json.NewEncoder(w).Encode(results)			
	}

	_ = json.NewDecoder(r.Body).Decode(&results)
	json.NewEncoder(w).Encode(results)	
}

func DeleteWorkflowResults(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var id string = params["id"]

	resultsDir := dataDir + "/" + profileName + "/" + configName + "/" + id

	var result util.Result
	err := util.RecursiveDirDelete(resultsDir)
	if err != nil {
		log.Println(err.Error())
		result = util.SetResultMessage(1,"ERROR",err.Error())
	} else {
		result = util.SetResultMessage(0,"INFO","Workflow results " + resultsDir + " deleted")
	}
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)	
}