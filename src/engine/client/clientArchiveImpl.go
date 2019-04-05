package client

import (
	"encoding/json"
	"engine/util"
	"log"
	"net/http"
	"bytes"
	"strings"
)

func ArchivePluginList(pluginType string,config util.Config) []string {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/pluginList/" + pluginType, b)
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

	var plugins []string

	if err := json.NewDecoder(resp.Body).Decode(&plugins); err != nil {
		log.Println(err)
	}

	return plugins

}

func ArchivePluginInfo(config util.Config, pluginName,pluginType string) (util.PluginInfoResult) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/pluginInfo/" + pluginName + "/" + pluginType, b)
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

	var pluginInfoResult util.PluginInfoResult
	if err := json.NewDecoder(resp.Body).Decode(&pluginInfoResult); err != nil {
		log.Println(err)
	}

	//unmarshall json response to plugin struct
	//var plugin util.Plugin
	//messages := strings.Join(result.Messages, "\n")
	//pluginByteArray := []byte(messages)

	//json.Unmarshal(pluginByteArray, &plugin)

	return pluginInfoResult
}

func Archive(config util.Config) util.Result {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/archive", b)
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

func ArchiveList(profileName,configName,policyName string,config util.Config) (util.ResultSimple, []util.Archive) {
	config = SetAdditionalConfigParams(profileName,configName,policyName,config)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/archiveList", b)
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

	var result util.ResultSimple
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	//unmarshall json response to plugin struct
	var archives []util.Archive
	messages := strings.Join(result.Messages, "\n")
	backupByteArray := []byte(messages)

	json.Unmarshal(backupByteArray, &archives)

	return result, archives
}

func ArchiveDelete(config util.Config) util.Result {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/archiveDelete", b)
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