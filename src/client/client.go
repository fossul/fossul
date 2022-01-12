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
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/fossul/fossul/src/engine/util"
)

type Auth struct {
	ServerHostname  string `json:"serverHostname,omitempty"`
	ServerPort      string `json:"serverPort,omitempty"`
	AppHostname     string `json:"appHostname,omitempty"`
	AppPort         string `json:"appPort,omitempty"`
	StorageHostname string `json:"storageHostname,omitempty"`
	StoragePort     string `json:"storagePort,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
}

func GetWorkflowStatus(auth Auth, profileName, configName string, id int) (util.WorkflowStatusResult, error) {
	var workflowStatusResult util.WorkflowStatusResult
	idToString := util.IntToString(id)
	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/getWorkflowStatus/"+profileName+"/"+configName+"/"+idToString, nil)
	if err != nil {
		return workflowStatusResult, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return workflowStatusResult, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&workflowStatusResult); err != nil {
			return workflowStatusResult, err
		}
	} else {
		return workflowStatusResult, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return workflowStatusResult, nil
}

func GetWorkflowStepResults(auth Auth, profileName, configName string, workflowId int, step int) ([]util.Result, error) {
	var results []util.Result
	w := util.IntToString(workflowId)
	s := util.IntToString(step)
	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/getWorkflowStepResults/"+profileName+"/"+configName+"/"+w+"/"+s, nil)
	if err != nil {
		return results, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return results, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
			return results, err
		}
	} else {
		return results, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return results, nil
}

func DeleteWorkflowResults(auth Auth, profileName, configName string, workflowId string) (util.Result, error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/deleteWorkflowResults/"+profileName+"/"+configName+"/"+workflowId, nil)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func GetServerServiceStatus(auth Auth) (util.Status, error) {
	var status util.Status

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/status", nil)
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return status, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			return status, err
		}
	} else {
		return status, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return status, nil
}

func GetAppServiceStatus(auth Auth) (util.Status, error) {
	var status util.Status

	req, err := http.NewRequest("GET", "http://"+auth.AppHostname+":"+auth.AppPort+"/status", nil)
	if err != nil {
		return status, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return status, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			return status, err
		}
	} else {
		return status, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return status, nil

}

func GetStorageServiceStatus(auth Auth) (util.Status, error) {
	var status util.Status

	req, err := http.NewRequest("GET", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/status", nil)
	if err != nil {
		return status, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return status, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			return status, err
		}
	} else {
		return status, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return status, nil

}

func StartBackupWorkflowLocalConfig(auth Auth, profileName, configName, policyName string, config util.Config) (util.WorkflowResult, error) {
	var result util.WorkflowResult
	config = SetAdditionalConfigParams(profileName, configName, policyName, config)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/startBackupWorkflowLocalConfig", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil

}

func StartBackupWorkflow(auth Auth, profileName, configName, policyName string) (util.WorkflowResult, error) {
	var result util.WorkflowResult

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/startBackupWorkflow/"+profileName+"/"+configName+"/"+policyName, nil)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil

}

func StartRestoreWorkflowLocalConfig(auth Auth, profileName, configName, policyName, selectedWorkflowId string, config util.Config) (util.WorkflowResult, error) {
	var result util.WorkflowResult
	config = SetAdditionalConfigParams(profileName, configName, policyName, config)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/startRestoreWorkflowLocalConfig/"+selectedWorkflowId, b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil

}

func StartRestoreWorkflow(auth Auth, profileName, configName, policyName, selectedWorkflowId string) (util.WorkflowResult, error) {
	var result util.WorkflowResult

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/startRestoreWorkflow/"+profileName+"/"+configName+"/"+policyName+"/"+selectedWorkflowId, nil)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil

}

func ServerBackupDelete(auth Auth, profileName, configName, policyName, selectedWorkflowId string) (util.Result, error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/deleteBackup/"+profileName+"/"+configName+"/"+policyName+"/"+selectedWorkflowId, nil)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil

}

func SetAdditionalConfigParams(profileName, configName, policyName string, config util.Config) util.Config {
	config.ProfileName = profileName
	config.ConfigName = configName

	backupRetention := util.GetBackupRetention(policyName, config.BackupRetentions)
	config.SelectedBackupRetention = backupRetention
	archiveRetention := util.GetArchiveRetention(policyName, config.ArchiveRetentions)
	config.SelectedArchiveRetention = archiveRetention
	config.SelectedBackupPolicy = policyName

	return config
}

func GetJobList(auth Auth, profileName, configName string) (util.Jobs, error) {

	var jobs util.Jobs

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/getJobs/"+profileName+"/"+configName, nil)
	if err != nil {
		return jobs, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return jobs, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
			return jobs, err
		}
	} else {
		return jobs, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return jobs, nil
}

func UpdateCustomBackupResource(auth Auth, profileName, configName, policyName, crName, workflowId string) (util.Result, error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/updateCustomBackupResource/"+profileName+"/"+configName+"/"+policyName+"/"+crName+"/"+workflowId, nil)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil

}
