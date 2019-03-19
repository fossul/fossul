package client

import (
	"encoding/json"
	"log"
	"engine/util"
	"net/http"
	"bytes"
)

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