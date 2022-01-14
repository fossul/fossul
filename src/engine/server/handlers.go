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
	"strings"

	"github.com/fossul/fossul/src/client"
	"github.com/fossul/fossul/src/client/k8s"
	"github.com/fossul/fossul/src/engine/util"
	"github.com/gorilla/mux"
)

// GetStatus godoc
// @Description Status and version information for the service
// @Accept  json
// @Produce  json
// @Success 200 {string} string
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /status [get]
func GetStatus(w http.ResponseWriter, r *http.Request) {
	var status util.Status
	status.Msg = "OK"
	status.Version = version

	json.NewEncoder(w).Encode(status)
}

// SendTrapSuccessCmd godoc
// @Description Execute command after successfull workflow execution
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /sendTrapSuccessCmd [post]
func SendTrapSuccessCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var messages []util.Message

	config, err := util.GetConfig(w, r)
	printConfigDebug(config)

	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't read config! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	if config.SendTrapSuccessCmd != "" {
		args := strings.Split(config.SendTrapSuccessCmd, ",")
		message := util.SetMessage("INFO", "Performing send trap success command ["+config.SendTrapSuccessCmd+"]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message, result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// SendTrapErrorCmd godoc
// @Description Execute command after failed workflow execution
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /sendTrapErrorCmd [post]
func SendTrapErrorCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var messages []util.Message

	config, err := util.GetConfig(w, r)
	printConfigDebug(config)

	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't read config! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	if config.SendTrapErrorCmd != "" {
		args := strings.Split(config.SendTrapErrorCmd, ",")
		message := util.SetMessage("INFO", "Performing send trap error command ["+config.SendTrapErrorCmd+"]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message, result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// GetJobs godoc
// @Description Get jobs (workflows) that have executed for a profile/config
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Jobs
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /getJobs/{profileName}/{configName} [get]
func GetJobs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]

	jobsDir := dataDir + "/" + profileName + "/" + configName

	var result util.Result
	var messages []util.Message

	var jobs util.Jobs
	var jobList []util.Job
	jobList, err := util.ListJobs(jobsDir)

	if err != nil {
		msg := util.SetMessage("ERROR", "Job list failed! "+err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		jobs.Result = result

		_ = json.NewDecoder(r.Body).Decode(&jobs)
		json.NewEncoder(w).Encode(jobs)
	}

	jobs.Result.Code = 0
	jobs.Jobs = jobList

	_ = json.NewDecoder(r.Body).Decode(&jobs)
	json.NewEncoder(w).Encode(jobs)
}

// DeleteBackup godoc
// @Description Delete a individual backup
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param policy path string true "name of backup policy"
// @Param workflowId path string true "workflow id"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /deleteBackup/{profileName}/{configName}/{policy}/{workflowId} [get]
func DeleteBackup(w http.ResponseWriter, r *http.Request) {
	auth := SetAuth()

	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var policyName string = params["policy"]
	var selectedWorkflowId string = params["workflowId"]

	var result util.Result
	var messages []util.Message

	config, err := util.GetConsolidatedConfig(configDir, profileName, configName, policyName)
	printConfigDebug(config)

	config.SelectedWorkflowId = util.StringToInt(selectedWorkflowId)

	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't read config! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	result, err = client.BackupDelete(auth, config)

	if err != nil {
		message := util.SetMessage("ERROR", err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	} else {
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// GetBackup godoc
// @Description Get a individual backup
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param workflowId path string true "workflow id"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.BackupByWorkflow
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /getBackup/{profileName}/{configName}/{workflowId} [get]
func GetBackup(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var selectedWorkflowId string = params["workflowId"]

	var backupByWorkflow util.BackupByWorkflow
	backup := &util.Backup{}
	var result util.Result
	var messages []util.Message

	backupFile := dataDir + "/" + profileName + "/" + configName + "/" + selectedWorkflowId + "/" + "backup"

	err := util.ReadGob(backupFile, &backup)
	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't read backup! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		backupByWorkflow.Result = result

		_ = json.NewDecoder(r.Body).Decode(&backupByWorkflow)
		json.NewEncoder(w).Encode(backupByWorkflow)

		return
	}

	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't read backup! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		backupByWorkflow.Result = result

		_ = json.NewDecoder(r.Body).Decode(&backupByWorkflow)
		json.NewEncoder(w).Encode(backupByWorkflow)

	} else {
		backupByWorkflow.Backup = *backup

		message := util.SetMessage("INFO", "Get backup for workflow id ["+selectedWorkflowId+"] completed successfully")
		messages = append(messages, message)

		result = util.SetResult(0, messages)
		backupByWorkflow.Result = result

		_ = json.NewDecoder(r.Body).Decode(&backupByWorkflow)
		json.NewEncoder(w).Encode(backupByWorkflow)
	}
}

// UpdateBackupCustomResource godoc
// @Description Update custom backup resource
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param policy path string true "name of backup policy"
// @Param crName path string true "name of custom resource"
// @Param op path string true "operation for patching resource"
// @Param specKey path string true "spec key we want to alter"
// @Param specValue path string true "spec key's value we want to alter"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /updateBackupCustomResource/{profileName}/{configName}/{policy}/{crName}/{op}/{specKey}/{specValue} [get]
func UpdateBackupCustomResource(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var policyName string = params["policy"]
	var crName string = params["crName"]
	var op string = params["op"]
	var specKey string = params["specKey"]
	var specValue string = params["specValue"]

	var result util.Result
	var messages []util.Message

	config, err := util.GetConsolidatedConfig(configDir, profileName, configName, policyName)
	printConfigDebug(config)

	msg := util.SetMessage("INFO", "Updating custom resource ["+crName+"] with workflowId ["+params["specValue"]+"]")
	messages = append(messages, msg)

	err = k8s.UpdateBackupCustomResource(config.AccessWithinCluster, profileName, crName, op, specKey, specValue)

	if err != nil {
		message := util.SetMessage("ERROR", err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	} else {
		msg := util.SetMessage("INFO", "Updating custom resource ["+crName+"] with workflowId ["+params["specValue"]+"] completed successfully")
		messages = append(messages, msg)

		result = util.SetResult(0, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}
