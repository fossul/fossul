package client

import (
	"encoding/json"
	"fossil/src/engine/util"
	"log"
	"net/http"
	"bytes"
	"strings"
	"errors"
)

func ArchivePluginList(auth Auth,pluginType string) ([]string,error) {
	var plugins []string

	req, err := http.NewRequest("POST", "http://" + auth.StorageHostname + ":" + auth.StoragePort + "/pluginList/" + pluginType, nil)
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

func ArchivePluginInfo(auth Auth,config util.Config, pluginName,pluginType string) (util.PluginInfoResult,error) {
	var pluginInfoResult util.PluginInfoResult
	
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://" + auth.StorageHostname + ":" + auth.StoragePort + "/pluginInfo/" + pluginName + "/" + pluginType, b)
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

func Archive(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result
	
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://" + auth.StorageHostname + ":" + auth.StoragePort + "/archive", b)
	if err != nil {
		log.Println("NewRequest: ", err)
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

func ArchiveList(auth Auth,profileName,configName,policyName string,config util.Config) (util.ResultSimple, []util.Archive, error) {
	var result util.ResultSimple
	var archives []util.Archive
	config = SetAdditionalConfigParams(profileName,configName,policyName,config)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://" + auth.StorageHostname + ":" + auth.StoragePort + "/archiveList", b)
	if err != nil {
		return result, archives, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, archives, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, archives, err
		}
	} else {
		return result,archives,errors.New("Http Status Error [" + resp.Status + "]")
	}

	//unmarshall json response to plugin struct
	messages := strings.Join(result.Messages, "\n")
	backupByteArray := []byte(messages)

	json.Unmarshal(backupByteArray, &archives)

	return result, archives, nil
}

func ArchiveDelete(auth Auth,config util.Config) (util.Result,error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://" + auth.StorageHostname + ":" + auth.StoragePort + "/archiveDelete", b)
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