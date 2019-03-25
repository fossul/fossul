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

func Archive(w http.ResponseWriter, r *http.Request) {
	var config util.Config = util.GetConfig(w,r)
	var plugin string = config.PluginDir + "/archive/" + config.ArchivePlugin

	if _, err := os.Stat(plugin); os.IsNotExist(err) {
		var errMsg string = "ERROR: Archive plugin does not exist"
		log.Println(err, errMsg)

		var messages []util.Message
		message := util.SetMessage("ERROR", errMsg + " " + err.Error())
		messages = append(messages, message)

		var result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}

	var result util.Result
	result = util.ExecutePlugin(config, "archive", plugin, "--action", "archive")	
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func ArchiveList(w http.ResponseWriter, r *http.Request) {
	var config util.Config = util.GetConfig(w,r)
	var plugin string = config.PluginDir + "/archive/" + config.ArchivePlugin

	if _, err := os.Stat(plugin); os.IsNotExist(err) {
		var errMsg string = "ERROR: Archive plugin does not exist"
		log.Println(err, errMsg)

		var messages []string
		message := fmt.Sprintf("ERROR %s %s",errMsg,err.Error())
		messages = append(messages, message)

		var result = util.SetResultSimple(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}

	var result util.ResultSimple
	result = util.ExecutePluginSimple(config, "archive", plugin, "--action", "archiveList")

	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func ArchiveDelete(w http.ResponseWriter, r *http.Request) {

	var config util.Config = util.GetConfig(w,r)
	var plugin string = config.PluginDir + "/archive/" + config.ArchivePlugin
	if _, err := os.Stat(plugin); os.IsNotExist(err) {
		var errMsg string = "ERROR: Archive plugin does not exist"
		log.Println(err, errMsg)

		var messages []util.Message
		message := util.SetMessage("ERROR", errMsg + " " + err.Error())
		messages = append(messages, message)

		var result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}

	var result util.Result
	result = util.ExecutePlugin(config, "archive", plugin, "--action", "archiveDelete")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func ArchiveCreateCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var config util.Config = util.GetConfig(w,r)

	if config.BackupCreateCmd != "" {
		args := strings.Split(config.ArchiveCreateCmd, ",")
		message := util.SetMessage("INFO", "Performing archive create command [" + config.ArchiveCreateCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

func ArchiveDeleteCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var config util.Config = util.GetConfig(w,r)

	if config.BackupDeleteCmd != "" {
		args := strings.Split(config.ArchiveDeleteCmd, ",")
		message := util.SetMessage("INFO", "Performing archive delete command [" +  config.ArchiveDeleteCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}