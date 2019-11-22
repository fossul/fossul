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
	"github.com/fossul/fossul/src/engine/util"
	"net/http"
)

func GetConfig(auth Auth, profileName, configName string) (util.ConfigResult, error) {
	var configResult util.ConfigResult

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/getConfig/"+profileName+"/"+configName, nil)
	if err != nil {
		return configResult, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return configResult, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&configResult); err != nil {
			return configResult, err
		}
	} else {
		return configResult, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return configResult, nil
}

func GetPluginConfig(auth Auth, profileName, configName, pluginName string) (util.ConfigMapResult, error) {
	var configMapResult util.ConfigMapResult

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/getPluginConfig/"+profileName+"/"+configName+"/"+pluginName, nil)
	if err != nil {
		return configMapResult, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return configMapResult, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&configMapResult); err != nil {
			return configMapResult, err
		}
	} else {
		return configMapResult, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return configMapResult, nil
}

func GetDefaultConfig(auth Auth) (util.ConfigResult, error) {
	var configResult util.ConfigResult

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/getDefaultConfig", nil)
	if err != nil {
		return configResult, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return configResult, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&configResult); err != nil {
			return configResult, err
		}
	} else {
		return configResult, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return configResult, nil
}

func GetDefaultPluginConfig(auth Auth, pluginName string) (util.ConfigMapResult, error) {
	var configMapResult util.ConfigMapResult

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/getDefaultPluginConfig/"+pluginName, nil)
	if err != nil {
		return configMapResult, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return configMapResult, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&configMapResult); err != nil {
			return configMapResult, err
		}
	} else {
		return configMapResult, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return configMapResult, nil
}

func AddConfig(auth Auth, profileName, configName string, config util.Config) (util.Result, error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	var result util.Result

	req, err := http.NewRequest("POST", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/addConfig/"+profileName+"/"+configName, b)
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

func AddPluginConfig(auth Auth, profileName, configName, pluginName string, configMap map[string]string) (util.Result, error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(configMap)

	var result util.Result

	req, err := http.NewRequest("POST", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/addPluginConfig/"+profileName+"/"+configName+"/"+pluginName, b)
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

func DeletePluginConfig(auth Auth, profileName, configName, pluginName string) (util.Result, error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/deletePluginConfig/"+profileName+"/"+configName+"/"+pluginName, nil)
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

func DeleteConfig(auth Auth, profileName, configName string) (util.Result, error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/deleteConfig/"+profileName+"/"+configName, nil)
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

func DeleteConfigDir(auth Auth, profileName, configName string) (util.Result, error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/deleteConfigDir/"+profileName+"/"+configName, nil)
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

func AddProfile(auth Auth, profileName string) (util.Result, error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/addProfile/"+profileName, nil)
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

func DeleteProfile(auth Auth, profileName string) (util.Result, error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/deleteProfile/"+profileName, nil)
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

func ListProfiles(auth Auth) (util.Result, error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/listProfiles", nil)
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

func ListConfigs(auth Auth, profileName string) (util.Result, error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/listConfigs/"+profileName, nil)
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

func ListPluginConfigs(auth Auth, profileName, configName string) (util.Result, error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/listPluginConfigs/"+profileName+"/"+configName, nil)
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
