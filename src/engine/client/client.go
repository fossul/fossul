package client

import (
	"encoding/json"
	"log"
	"engine/util"
	"net/http"
	"bytes"
)

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

func StartBackupWorkflow(profileName,configName,policyName string,config util.Config) (result []util.Result) {

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