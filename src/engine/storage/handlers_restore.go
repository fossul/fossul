package main

import (
	"encoding/json"
	"fossil/src/engine/util"
	"net/http"
	"os"
	"strings"
)

func Restore(w http.ResponseWriter, r *http.Request) {
	config,_ := util.GetConfig(w,r)
	pluginPath := util.GetPluginPath(config.StoragePlugin)
	var result util.Result
	var messages []util.Message

	if pluginPath == "" {
		var plugin string = pluginDir + "/storage/" + config.StoragePlugin

		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			var errMsg string = "Storage plugin does not exist"

			message := util.SetMessage("ERROR", errMsg + " " + err.Error())
			messages = append(messages, message)

			result = util.SetResult(1, messages)
			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}
		result = util.ExecutePlugin(config, "storage", plugin, "--action", "restore")	
		_ 	= json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		plugin,err := util.GetStorageInterface(pluginPath)
		if err != nil {
			message := util.SetMessage("ERROR", err.Error())
			messages = append(messages, message)

			var result = util.SetResult(1, messages)			
			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)		
		} else {
			setEnvResult := plugin.SetEnv(config)
			if setEnvResult.Code != 0 {
				_ = json.NewDecoder(r.Body).Decode(&setEnvResult)
				json.NewEncoder(w).Encode(setEnvResult)
			} else {
				result = plugin.Restore(config)
				messages = util.PrependMessages(setEnvResult.Messages,result.Messages)
				result.Messages = messages
	
				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)			
			}	
		}	
	}	
}

func RestoreCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	config,_ := util.GetConfig(w,r)

	if config.RestoreCmd != "" {
		args := strings.Split(config.RestoreCmd, ",")
		message := util.SetMessage("INFO", "Performing restore command [" + config.RestoreCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}