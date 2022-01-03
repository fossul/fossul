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
	"net/http"

	"github.com/fossul/fossul/src/engine/util"
)

func StoragePluginList(auth Auth, pluginType string) ([]string, error) {
	var plugins []string

	req, err := http.NewRequest("GET", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/pluginList/"+pluginType, nil)
	if err != nil {
		return plugins, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return plugins, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&plugins); err != nil {
			return plugins, err
		}
	} else {
		return plugins, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return plugins, nil

}

func StoragePluginInfo(auth Auth, config util.Config, pluginName, pluginType string) (util.PluginInfoResult, error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	var pluginInfoResult util.PluginInfoResult

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/pluginInfo/"+pluginName+"/"+pluginType, b)
	if err != nil {
		return pluginInfoResult, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return pluginInfoResult, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&pluginInfoResult); err != nil {
			return pluginInfoResult, err
		}
	} else {
		return pluginInfoResult, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return pluginInfoResult, nil
}

func Backup(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/backup", b)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return result, err
	}

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

func Restore(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/restore", b)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return result, err
	}

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

func Mount(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/mount", b)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return result, err
	}

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

func Unmount(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/unmount", b)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return result, err
	}

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

func BackupList(auth Auth, profileName, configName, policyName string, config util.Config) (util.Backups, error) {
	var backups util.Backups

	config = SetAdditionalConfigParams(profileName, configName, policyName, config)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/backupList", b)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return backups, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return backups, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&backups); err != nil {
			return backups, err
		}
	} else {
		return backups, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return backups, nil
}

func BackupDeleteWorkflow(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/backupDeleteWorkflow", b)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return result, err
	}

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

func BackupDelete(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/backupDelete", b)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return result, err
	}

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
