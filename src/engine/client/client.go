package client

import (
	"encoding/json"
	"log"
	"fossil/src/engine/util"
	"net/http"
	"bytes"
)

func GetWorkflowStatus(profileName,configName string,id int) (workflow util.Workflow) {

	idToString := util.IntToString(id)
	req, err := http.NewRequest("GET", "http://fossil-workflow:8000/getWorkflowStatus/" + profileName + "/" + configName + "/" + idToString, nil)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&workflow); err != nil {
		log.Println(err)
	}

	return workflow
}

func GetWorkflowStepResults(profileName,configName string,workflowId int, step int) (results []util.Result) {

	w := util.IntToString(workflowId)
	s := util.IntToString(step)
	req, err := http.NewRequest("GET", "http://fossil-workflow:8000/getWorkflowStepResults/" + profileName + "/" + configName + "/" + w + "/" + s, nil)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		log.Println(err)
	}

	return results
}

func DeleteWorkflowResults(profileName,configName string,workflowId string) util.Result {

	req, err := http.NewRequest("POST", "http://fossil-workflow:8000/deleteWorkflowResults/" + profileName + "/" + configName + "/" + workflowId, nil)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	var result util.Result

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	return result
}

func GetConfig(profileName,configName string) util.Config {

	req, err := http.NewRequest("GET", "http://fossil-workflow:8000/getConfig/" + profileName + "/" + configName, nil)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	var config util.Config
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		log.Println(err)
	}

	return config
}

func GetDefaultConfig() util.Config {

	req, err := http.NewRequest("GET", "http://fossil-workflow:8000/getDefaultConfig", nil)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	var config util.Config
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		log.Println(err)
	}

	return config
}

func GetDefaultPluginConfig(pluginName string) map[string]string {

	req, err := http.NewRequest("GET", "http://fossil-workflow:8000/getDefaultPluginConfig/" + pluginName, nil)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	var configMap map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&configMap); err != nil {
		log.Println(err)
	}

	return configMap
}

func GetWorkflowServiceStatus() util.Status {

	req, err := http.NewRequest("GET", "http://fossil-workflow:8000/status", nil)
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	var status util.Status

	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		log.Println(err)
	}

	return status

}

func GetAppServiceStatus() util.Status {

	req, err := http.NewRequest("GET", "http://fossil-app:8001/status", nil)
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	var status util.Status

	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		log.Println(err)
	}

	return status

}

func GetStorageServiceStatus() util.Status {

	req, err := http.NewRequest("GET", "http://fossil-storage:8002/status", nil)
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	var status util.Status

	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		log.Println(err)
	}

	return status

}

func StartBackupWorkflow(profileName,configName,policyName string,config util.Config) (result util.WorkflowResult) {

	config = SetAdditionalConfigParams(profileName,configName,policyName,config)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-workflow:8000/startBackupWorkflow", b)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	return result

}

func SetAdditionalConfigParams(profileName, configName, policyName string, config util.Config) util.Config {
	config.ProfileName = profileName
	config.ConfigName = configName

	backupRetention := util.GetBackupRetention(policyName,config.BackupRetentions)
	config.SelectedBackupRetention = backupRetention
	config.SelectedBackupPolicy = policyName

	return config
}