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
	"encoding/json"
	"net/http"
	"time"

	"github.com/fossul/fossul/src/engine/util"
	"github.com/gorilla/mux"
)

// StartBackupWorkflowLocalConfig godoc
// @Description Start backup workflow using local config
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.WorkflowResult
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /startBackupWorkflowLocalConfig [post]
func StartBackupWorkflowLocalConfig(w http.ResponseWriter, r *http.Request) {
	var workflowResult util.WorkflowResult
	workflow := &util.Workflow{}
	workflow.Id = util.GetWorkflowId()
	workflow.Type = "backup"
	workflow.Status = "RUNNING"

	var timestamp string = time.Now().Format(time.RFC3339)
	workflow.Timestamp = timestamp

	workflowResult.Id = workflow.Id

	var result util.Result
	var messages []util.Message

	config, err := util.GetConfig(w, r)
	printConfigDebug(config)

	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't read config! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		workflowResult.Result = result

		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)

		return
	}

	config.WorkflowId = util.IntToString(workflow.Id)
	config.WorkflowType = workflow.Type
	config.WorkflowTimestamp = util.GetTimestamp()
	workflow.Policy = config.SelectedBackupPolicy

	_, ok := runningWorkflowMap.Load(config.ProfileName + "-" + config.ConfigName)
	if ok {
		result = util.SetResultMessage(1, "ERROR", "Workflow id ["+util.IntToString(workflow.Id)+"] failed to start. Another workflow is running under profile ["+config.ProfileName+"] config ["+config.ConfigName+"]")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)
	} else {
		runningWorkflowMap.Store(config.ProfileName+"-"+config.ConfigName, config.SelectedBackupPolicy)

		go func() {
			startBackupWorkflowImpl(dataDir, config, workflow)
		}()

		result = util.SetResultMessage(0, "INFO", "Workflow id ["+util.IntToString(workflow.Id)+"] started successfully")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)
	}
}

// StartBackupWorkflow godoc
// @Description Start backup workflow using local config
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param policy path string true "name of backup policy"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.WorkflowResult
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /startBackupWorkflow/{profileName}/{configName}/{policy} [post]
func StartBackupWorkflow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var policyName string = params["policy"]

	var workflowResult util.WorkflowResult
	workflow := &util.Workflow{}
	workflow.Id = util.GetWorkflowId()
	workflow.Type = "backup"
	workflow.Policy = policyName
	workflow.Status = "RUNNING"

	var timestamp string = time.Now().Format(time.RFC3339)
	workflow.Timestamp = timestamp

	workflowResult.Id = workflow.Id

	config, err := util.GetConsolidatedConfig(configDir, profileName, configName, policyName)
	printConfigDebug(config)

	if err != nil {
		result := util.SetResultMessage(1, "ERROR", "Workflow id ["+util.IntToString(workflow.Id)+"] failed to start. Couldn't read config using profile ["+profileName+"] config ["+configName+"] "+err.Error())
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)
	}

	config.WorkflowId = util.IntToString(workflow.Id)
	config.WorkflowType = workflow.Type
	config.WorkflowTimestamp = util.GetTimestamp()

	_, ok := runningWorkflowMap.Load(config.ProfileName + "-" + config.ConfigName)
	if ok {
		result := util.SetResultMessage(1, "ERROR", "Workflow id ["+util.IntToString(workflow.Id)+"] failed to start. Another workflow is running under profile ["+config.ProfileName+"] config ["+config.ConfigName+"]")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)
	} else {
		runningWorkflowMap.Store(config.ProfileName+"-"+config.ConfigName, config.SelectedBackupPolicy)

		go func() {
			startBackupWorkflowImpl(dataDir, config, workflow)
		}()

		result := util.SetResultMessage(0, "INFO", "Workflow id ["+util.IntToString(workflow.Id)+"] started successfully")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)
	}
}

// StartRestoreWorkflowLocalConfig godoc
// @Description Start restore workflow using local config
// @Param workflowId body workflowId true int
// @Accept  json
// @Produce  json
// @Success 200 {object} util.WorkflowResult
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /startRestoreWorkflowLocalConfig [post]
func StartRestoreWorkflowLocalConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var selectedWorkflowId string = params["workflowId"]

	var workflowResult util.WorkflowResult
	workflow := &util.Workflow{}
	workflow.Id = util.GetWorkflowId()
	workflow.Type = "restore"
	workflow.Status = "RUNNING"

	var timestamp string = time.Now().Format(time.RFC3339)
	workflow.Timestamp = timestamp

	workflowResult.Id = workflow.Id

	var result util.Result
	var messages []util.Message

	config, err := util.GetConfig(w, r)
	printConfigDebug(config)

	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't read config! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		workflowResult.Result = result

		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)

		return
	}

	config.WorkflowId = util.IntToString(workflow.Id)
	config.WorkflowType = workflow.Type
	config.WorkflowTimestamp = util.GetTimestamp()
	config.SelectedWorkflowId = util.StringToInt(selectedWorkflowId)
	workflow.Policy = config.SelectedBackupPolicy

	_, ok := runningWorkflowMap.Load(config.ProfileName + "-" + config.ConfigName)
	if ok {
		result = util.SetResultMessage(1, "ERROR", "Workflow id ["+util.IntToString(workflow.Id)+"] failed to start. Another workflow is running under profile ["+config.ProfileName+"] config ["+config.ConfigName+"]")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)
	} else {
		runningWorkflowMap.Store(config.ProfileName+"-"+config.ConfigName, config.SelectedBackupPolicy)

		go func() {
			startRestoreWorkflowImpl(dataDir, config, workflow)
		}()

		result = util.SetResultMessage(0, "INFO", "Workflow id ["+util.IntToString(workflow.Id)+"] started successfully")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)
	}
}

// StartRestoreWorkflow godoc
// @Description Start restore workflow using local config
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param policy path string true "name of backup policy"
// @Param workflowId path string true "workflow id"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.WorkflowResult
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /startRestoreWorkflow/{profileName}/{configName}/{policy}/{workflowId} [post]
func StartRestoreWorkflow(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var policyName string = params["policy"]
	var selectedWorkflowId string = params["workflowId"]

	var workflowResult util.WorkflowResult
	workflow := &util.Workflow{}
	workflow.Id = util.GetWorkflowId()
	workflow.Type = "restore"
	workflow.Policy = policyName
	workflow.Status = "RUNNING"

	var timestamp string = time.Now().Format(time.RFC3339)
	workflow.Timestamp = timestamp

	workflowResult.Id = workflow.Id

	config, err := util.GetConsolidatedConfig(configDir, profileName, configName, policyName)
	printConfigDebug(config)

	if err != nil {
		result := util.SetResultMessage(1, "ERROR", "Workflow id ["+util.IntToString(workflow.Id)+"] failed to start. Couldn't read config using profile ["+profileName+"] config ["+configName+"]")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)
	}

	config.WorkflowId = util.IntToString(workflow.Id)
	config.WorkflowType = workflow.Type
	config.WorkflowTimestamp = util.GetTimestamp()
	config.SelectedWorkflowId = util.StringToInt(selectedWorkflowId)

	_, ok := runningWorkflowMap.Load(config.ProfileName + "-" + config.ConfigName)
	if ok {
		result := util.SetResultMessage(1, "ERROR", "Workflow id ["+util.IntToString(workflow.Id)+"] failed to start. Another workflow is running under profile ["+config.ProfileName+"] config ["+config.ConfigName+"]")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)
	} else {
		runningWorkflowMap.Store(config.ProfileName+"-"+config.ConfigName, config.SelectedBackupPolicy)

		go func() {
			startRestoreWorkflowImpl(dataDir, config, workflow)
		}()

		result := util.SetResultMessage(0, "INFO", "Workflow id ["+util.IntToString(workflow.Id)+"] started successfully")
		workflowResult.Result = result
		_ = json.NewDecoder(r.Body).Decode(&workflowResult)
		json.NewEncoder(w).Encode(workflowResult)
	}
}

// GetWorkflowStatus godoc
// @Description Get workflow status
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param id path string true "workflow id"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.WorkflowStatusResult
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /getWorkflowStatus/{profileName}/{configName}/{id} [get]
func GetWorkflowStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var id string = params["id"]

	var workflowStatusResult util.WorkflowStatusResult

	resultsDir := dataDir + "/" + profileName + "/" + configName + "/" + id

	workflow := &util.Workflow{}

	workflowFile := resultsDir + "/workflow"
	err := util.ReadGob(workflowFile, &workflow)
	if err != nil {
		result := util.SetResultMessage(1, "ERROR", "Couldn't get status for workflow id ["+id+"] "+err.Error())
		workflowStatusResult.Result = result

		_ = json.NewDecoder(r.Body).Decode(&workflowStatusResult)
		json.NewEncoder(w).Encode(workflowStatusResult)
	} else {
		workflowStatusResult.Workflow = *workflow

		_ = json.NewDecoder(r.Body).Decode(&workflowStatusResult)
		json.NewEncoder(w).Encode(workflowStatusResult)
	}
}

// GetWorkflowStepResults godoc
// @Description Get workflow step results
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param workflowId path string true "workflow id"
// @Param stepId path string true "step id"
// @Accept  json
// @Produce  json
// @Success 200 {array} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /getWorkflowStepResults/{profileName}/{configName}/{workflowId}/{stepId} [get]
func GetWorkflowStepResults(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var workflowId string = params["workflowId"]
	var stepId string = params["stepId"]

	var results []util.Result
	resultsFile := dataDir + "/" + profileName + "/" + configName + "/" + workflowId + "/" + stepId

	var result util.Result
	err := util.ReadGob(resultsFile, &result)
	results = append(results, result)
	if err != nil {
		var results []util.Result
		errorResult := util.SetResultMessage(1, "ERROR", err.Error())
		results = append(results, errorResult)

		_ = json.NewDecoder(r.Body).Decode(&results)
		json.NewEncoder(w).Encode(results)
	}

	_ = json.NewDecoder(r.Body).Decode(&results)
	json.NewEncoder(w).Encode(results)
}

// DeleteWorkflowResults godoc
// @Description Delete workflow results for profile/config
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param id path string true "workflow id"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /deleteWorkflowResults/{profileName}/{configName}/{id} [get]
func DeleteWorkflowResults(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var id string = params["id"]

	resultsDir := dataDir + "/" + profileName + "/" + configName + "/" + id

	var result util.Result
	err := util.RecursiveDirDelete(resultsDir)
	if err != nil {
		result = util.SetResultMessage(1, "ERROR", err.Error())
	} else {
		result = util.SetResultMessage(0, "INFO", "Workflow results "+resultsDir+" deleted")
	}
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}
