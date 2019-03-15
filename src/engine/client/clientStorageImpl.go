package client

import (
	"encoding/json"
	"engine/util"
	"log"
	"net/http"
	"bytes"
	"strings"
)

func StoragePluginList(config util.Config) []string {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/pluginList", b)
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

func StoragePluginInfo(config util.Config, pluginName string) (util.ResultSimple, util.Plugin) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/pluginInfo/" + pluginName, b)
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
	var plugin util.Plugin
	messages := strings.Join(result.Messages, "\n")
	pluginByteArray := []byte(messages)

	json.Unmarshal(pluginByteArray, &plugin)

	return result, plugin
}

func Backup(config util.Config) util.Result {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/backup", b)
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

func BackupList(profileName,configName string,config util.Config) (util.ResultSimple, []util.Backup) {
	config = SetAdditionalConfigParams(profileName,configName,config)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/backupList", b)
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
	var backups []util.Backup
	messages := strings.Join(result.Messages, "\n")
	backupByteArray := []byte(messages)

	json.Unmarshal(backupByteArray, &backups)

	return result, backups
}


/*
func BackupList(profileName,configName string,config util.Config) util.ResultSimple {
	config = setAdditionalConfigParams(profileName,configName,config)

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/backupList", b)
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

	return result
}

*/

func BackupDelete(config util.Config) util.Result {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-storage:8002/backupDelete", b)
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