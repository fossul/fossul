package client

import (
	"encoding/json"
	"fossil/src/engine/util"
	"net/http"
	"bytes"
	"errors"
)

func GetConfig(auth Auth,profileName,configName string) (util.Config,error) {
	var config util.Config

	req, err := http.NewRequest("GET", "http://fossil-server:8000/getConfig/" + profileName + "/" + configName, nil)
	if err != nil {
		return config,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return config,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
			return config,err
		}
	} else {
		return config,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return config,nil
}

func GetPluginConfig(auth Auth,profileName,configName,pluginName string) (map[string]string,error) {
	var configMap map[string]string

	req, err := http.NewRequest("GET", "http://fossil-server:8000/getPluginConfig/" + profileName + "/" + configName + "/" + pluginName, nil)
	if err != nil {
		return configMap,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return configMap,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&configMap); err != nil {
			return configMap,err
		}
	} else {
		return configMap,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return configMap,nil
}

func GetDefaultConfig(auth Auth) (util.Config,error) {
	var config util.Config

	req, err := http.NewRequest("GET", "http://fossil-server:8000/getDefaultConfig", nil)
	if err != nil {
		return config,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return config,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
			return config,err
		}
	} else {
		return config,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return config,nil
}

func GetDefaultPluginConfig(auth Auth,pluginName string) (map[string]string,error) {
	var configMap map[string]string

	req, err := http.NewRequest("GET", "http://fossil-server:8000/getDefaultPluginConfig/" + pluginName, nil)
	if err != nil {
		return configMap,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return configMap,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&configMap); err != nil {
			return configMap,err
		}
	} else {
		return configMap,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return configMap,nil
}

func AddConfig(auth Auth,profileName,configName string,config util.Config) (util.Result,error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	var result util.Result

	req, err := http.NewRequest("GET", "http://fossil-server:8000/addConfig/" + profileName + "/" + configName, b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func AddPluginConfig(auth Auth,profileName,configName,pluginName string,configMap map[string]string) (util.Result,error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(configMap)

	var result util.Result

	req, err := http.NewRequest("GET", "http://fossil-server:8000/addPluginConfig/" + profileName + "/" + configName + "/" + pluginName, b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func DeleteConfig(auth Auth,profileName,configName string) (util.Result,error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://fossil-server:8000/deleteConfig/" + profileName + "/" + configName, nil)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func AddProfile(auth Auth,profileName string) (util.Result,error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://fossil-server:8000/addProfile/" + profileName, nil)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func DeleteProfile(auth Auth,profileName string) (util.Result,error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://fossil-server:8000/deleteProfile/" + profileName, nil)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func ListProfiles(auth Auth) (util.Result,error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://fossil-server:8000/listProfiles", nil)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func ListConfigs(auth Auth,profileName string) (util.Result,error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://fossil-server:8000/listConfigs/" + profileName, nil)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func ListPluginConfigs(auth Auth,profileName,configName string) (util.Result,error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://fossil-server:8000/listPluginConfigs/" + profileName + "/" + configName, nil)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}