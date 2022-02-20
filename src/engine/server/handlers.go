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
	"os"
	"strings"

	"fmt"

	"github.com/fossul/fossul/src/client"
	"github.com/fossul/fossul/src/client/k8s"
	"github.com/fossul/fossul/src/engine/util"
	"github.com/fossul/fossul/src/plugins/pluginUtil"
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

	backupFile := dataDir + "/" + profileName + "/" + configName + "/" + selectedWorkflowId + "/backup"

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

		err = os.Remove(backupFile)
		if err != nil {
			message := util.SetMessage("ERROR", err.Error())
			messages = append(messages, message)

			result = util.SetResult(1, messages)

			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}

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

// CreateBackupCustomResource godoc
// @Description Update custom backup resource
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param policy path string true "name of backup policy"
// @Param workflowId path string true "workflow id"
// @Param timestamp path string true "timestamp of custom resource creation"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /createBackupCustomResource/{profileName}/{configName}/{policy}/{workflowId}/{timestamp} [get]
func CreateBackupCustomResource(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var policyName string = params["policy"]
	var workflowId string = params["workflowId"]
	var timestamp string = params["timestamp"]

	var result util.Result
	var messages []util.Message

	config, err := util.GetConsolidatedConfig(configDir, profileName, configName, policyName)
	printConfigDebug(config)

	backupName := util.GetBackupName(config.StoragePluginParameters["BackupName"], config.SelectedBackupPolicy, workflowId, timestamp)

	msg := util.SetMessage("INFO", "Creating backup custom resource ["+backupName+"]")
	messages = append(messages, msg)

	err = k8s.CreateBackupCustomResource(config.AccessWithinCluster, profileName, backupName, profileName, configName, policyName)

	if err != nil {
		message := util.SetMessage("ERROR", err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	} else {
		msg := util.SetMessage("INFO", "Creating backup custom resource ["+backupName+"] completed successfully")
		messages = append(messages, msg)

		result = util.SetResult(0, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// DeleteBackupCustomResource godoc
// @Description Update custom backup resource
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param policy path string true "name of backup policy"
// @Param crName path string true "name of custom resource"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /deleteBackupCustomResource/{profileName}/{configName}/{policy}/{crName} [get]
func DeleteBackupCustomResource(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var policyName string = params["policy"]
	var crName string = params["crName"]

	var result util.Result
	var messages []util.Message

	config, err := util.GetConsolidatedConfig(configDir, profileName, configName, policyName)
	printConfigDebug(config)

	msg := util.SetMessage("INFO", "Deleting backup custom resource ["+crName+"]")
	messages = append(messages, msg)

	err = k8s.DeleteBackupCustomResource(config.AccessWithinCluster, profileName, crName)

	if err != nil {
		message := util.SetMessage("ERROR", err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	} else {
		msg := util.SetMessage("INFO", "Deleting backup custom resource ["+crName+"] completed successfully")
		messages = append(messages, msg)

		result = util.SetResult(0, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// BackupCustomResourceRetention godoc
// @Description Update custom backup resource
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param policy path string true "name of backup policy"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /deleteBackupCustomResource/{profileName}/{configName}/{policy} [get]
func BackupCustomResourceRetention(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var policyName string = params["policy"]

	var result util.Result
	var messages []util.Message

	config, err := util.GetConsolidatedConfig(configDir, profileName, configName, policyName)
	printConfigDebug(config)

	var backupList []string
	backupCustomResourceList, err := k8s.ListBackupCustomResources(config.AccessWithinCluster, profileName)
	for _, backupCustomResource := range backupCustomResourceList.Items {
		backupCustomResourceName := backupCustomResource.GetName()

		backupList = append(backupList, backupCustomResourceName)
	}

	backups, err := pluginUtil.ListCustomResourceBackups(backupList, config.StoragePluginParameters["BackupName"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	backupsByPolicy := util.GetBackupsByPolicy(config.SelectedBackupPolicy, backups)
	backupCount := len(backupsByPolicy)

	if backupCount > config.SelectedBackupRetention {
		count := 1
		for backup := range pluginUtil.ReverseBackupList(backupsByPolicy) {
			if count > config.SelectedBackupRetention {
				msg := util.SetMessage("INFO", fmt.Sprintf("Number of backups [%d] greater than backup retention [%d]", backupCount, config.SelectedBackupRetention))
				messages = append(messages, msg)
				backupCount = backupCount - 1

				for _, content := range backup.Contents {
					backupName := content.Data
					msg = util.SetMessage("INFO", "Deleting backup "+backupName)
					messages = append(messages, msg)

					err := k8s.DeleteBackupCustomResource(config.AccessWithinCluster, config.ProfileName, backupName)
					if err != nil {
						msg := util.SetMessage("ERROR", err.Error())
						messages = append(messages, msg)
						result = util.SetResult(1, messages)

						_ = json.NewDecoder(r.Body).Decode(&result)
						json.NewEncoder(w).Encode(result)

						return
					}

					msg = util.SetMessage("INFO", "Backup "+backupName+" deleted successfully")
					messages = append(messages, msg)
				}
			}
			count = count + 1
		}
	} else {
		msg := util.SetMessage("INFO", fmt.Sprintf("Backup deletion skipped, there are [%d] backups but backup retention is [%d]", backupCount, config.SelectedBackupRetention))
		messages = append(messages, msg)
	}

	result = util.SetResult(0, messages)

	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}
