package main

import (
	"encoding/json"
	"engine/util"
	"net/http"
	"log"
	"os"
	"fmt"
	"strings"
)

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
		message := util.SetMessage("INFO", "Performing backup create command [" + config.BackupCreateCmd + "]")

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
		message := util.SetMessage("INFO", "Performing backup delete command [" + config.BackupDeleteCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}