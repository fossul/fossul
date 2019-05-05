package client

import (
	"encoding/json"
	"fossil/src/engine/util"
	"net/http"
	"bytes"
	"errors"
)

func StoragePluginList(auth Auth,pluginType string,config util.Config) ([]string,error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	var plugins []string

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/pluginList/" + pluginType, b)
	if err != nil {
		return plugins,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return plugins,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&plugins); err != nil {
			return plugins,err
		}
	} else {
		return plugins,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return plugins,nil

}

func StoragePluginInfo(auth Auth,config util.Config, pluginName,pluginType string) (util.PluginInfoResult,error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	var pluginInfoResult util.PluginInfoResult

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/pluginInfo/" + pluginName + "/" + pluginType, b)
	if err != nil {
		return pluginInfoResult,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return pluginInfoResult,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&pluginInfoResult); err != nil {
			return pluginInfoResult,err
		}
	} else {
		return pluginInfoResult,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return pluginInfoResult,nil
}

func Backup(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/backup", b)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return result,err
	}

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

func BackupList(auth Auth,profileName,configName,policyName string,config util.Config) (util.Backups,error) {
	var backups util.Backups

	config = SetAdditionalConfigParams(profileName,configName,policyName,config)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/backupList", b)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return backups,err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return backups,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&backups); err != nil {
			return backups,err
		}
	} else {
		return backups,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return backups,nil
}

func BackupDelete(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/backupDelete", b)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return result,err
	}

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