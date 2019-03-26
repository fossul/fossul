package main

import (
	"encoding/json"
	"engine/util"
	"net/http"
	"log"
	"os"
	"strings"
)

func UnquiesceCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var config util.Config = util.GetConfig(w,r)

	if config.PreAppQuiesceCmd != "" {
		args := strings.Split(config.AppUnquiesceCmd, ",")
		message := util.SetMessage("INFO", "Performing unquiesce command [" + config.PreAppQuiesceCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

func PreUnquiesceCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var config util.Config = util.GetConfig(w,r)

	if config.PreAppQuiesceCmd != "" {
		args := strings.Split(config.PreAppUnquiesceCmd, ",")
		message := util.SetMessage("INFO", "Performing pre unquiesce command [" + config.PreAppQuiesceCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

func Unquiesce(w http.ResponseWriter, r *http.Request) {

	var config util.Config = util.GetConfig(w,r)
	var plugin string = config.PluginDir + "/app/" + config.AppPlugin
	
	if _, err := os.Stat(plugin); os.IsNotExist(err) {
		var errMsg string = "ERROR: App plugin does not exist"
		log.Println(err, errMsg)

		var messages []util.Message
		message := util.SetMessage("ERROR", errMsg + " " + err.Error())
		messages = append(messages, message)

		var result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}

	var result util.Result
	result = util.ExecutePlugin(config, "app", plugin, "--action", "unquiesce")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func PostUnquiesceCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var config util.Config = util.GetConfig(w,r)

	if config.PreAppQuiesceCmd != "" {
		args := strings.Split(config.PostAppUnquiesceCmd, ",")
		message := util.SetMessage("INFO", "Performing post unquiesce command [" + config.PreAppQuiesceCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}
