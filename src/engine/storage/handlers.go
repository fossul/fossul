package main

import (
	"encoding/json"
	"engine/util"
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"os"
	"io/ioutil"
	"fmt"
	"strings"
)

func GetStatus(w http.ResponseWriter, r *http.Request) {
	var status = util.Status{Msg: "OK"}
	json.NewEncoder(w).Encode(status)
}

func PluginList(w http.ResponseWriter, r *http.Request) {
	var config util.Config = util.GetConfig(w,r)
	pluginDir := config.PluginDir + "/storage"

	var plugins []string
	fileInfo, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		log.Println(err)
	}

	for _, file := range fileInfo {
		plugins = append(plugins, file.Name())
	}

	_ = json.NewDecoder(r.Body).Decode(&plugins)
	json.NewEncoder(w).Encode(plugins)
}

func PluginInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var pluginName string = params["plugin"]

	var config util.Config = util.GetConfig(w,r)
	var plugin string = config.PluginDir + "/storage/" + pluginName

	var result util.ResultSimple
	result = util.ExecutePluginSimple(config, "storage,", plugin, "--action", "info")

	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func Backup(w http.ResponseWriter, r *http.Request) {
	var config util.Config = util.GetConfig(w,r)
	var plugin string = config.PluginDir + "/storage/" + config.StoragePlugin

	if _, err := os.Stat(plugin); os.IsNotExist(err) {
		var errMsg string = "ERROR: Storage plugin does not exist"
		log.Println(err, errMsg)

		var messages []util.Message
		message := util.SetMessage("ERROR", errMsg + " " + err.Error())
		messages = append(messages, message)

		var result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}

	var result util.Result
	result = util.ExecutePlugin(config, "storage", plugin, "--action", "backup")	
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func BackupList(w http.ResponseWriter, r *http.Request) {
	var config util.Config = util.GetConfig(w,r)
	var plugin string = config.PluginDir + "/storage/" + config.StoragePlugin

	if _, err := os.Stat(plugin); os.IsNotExist(err) {
		var errMsg string = "ERROR: Storage plugin does not exist"
		log.Println(err, errMsg)

		var messages []string
		message := fmt.Sprintf("ERROR %s %s",errMsg,err.Error())
		messages = append(messages, message)

		var result = util.SetResultSimple(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}

	var result util.ResultSimple
	result = util.ExecutePluginSimple(config, "storage", plugin, "--action", "backupList")

	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func BackupDelete(w http.ResponseWriter, r *http.Request) {

	var config util.Config = util.GetConfig(w,r)
	var plugin string = config.PluginDir + "/storage/" + config.StoragePlugin
	if _, err := os.Stat(plugin); os.IsNotExist(err) {
		var errMsg string = "ERROR: Storage plugin does not exist"
		log.Println(err, errMsg)

		var messages []util.Message
		message := util.SetMessage("ERROR", errMsg + " " + err.Error())
		messages = append(messages, message)

		var result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}

	var result util.Result
	result = util.ExecutePlugin(config, "storage", plugin, "--action", "backupDelete")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func BackupCreateCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var config util.Config = util.GetConfig(w,r)

	if config.BackupCreateCmd != "" {
		args := strings.Split(config.BackupCreateCmd, ",")
		message := util.SetMessage("INFO", "Performing backup create command")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

func BackupDeleteCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var config util.Config = util.GetConfig(w,r)

	if config.BackupDeleteCmd != "" {
		args := strings.Split(config.BackupDeleteCmd, ",")
		message := util.SetMessage("INFO", "Performing backup delete command " + config.BackupDeleteCmd)

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}
