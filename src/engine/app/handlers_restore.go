package main

import (
	"encoding/json"
	"fossil/src/engine/util"
	"net/http"
	"strings"
	"os"
)

func PreRestore(w http.ResponseWriter, r *http.Request) {

	config,_ := util.GetConfig(w,r)
	printConfigDebug(config)

	pluginPath := util.GetPluginPath(config.AppPlugin)
	var messages []util.Message

	if pluginPath == "" {
		var plugin string = pluginDir + "/app/" + config.AppPlugin
		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			var errMsg string = "\nApp plugin does not exist: " + plugin
	
			message := util.SetMessage("ERROR", errMsg + " " + err.Error())
			messages = append(messages, message)
	
			var result = util.SetResult(1, messages)
	
			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}
	
		var result util.Result
		result = util.ExecutePlugin(config, "app", plugin, "--action", "preRestore")
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {	
		var result util.Result
		plugin,err := util.GetAppInterface(pluginPath)
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
				result = plugin.PreRestore(config)
				messages = util.PrependMessages(setEnvResult.Messages,result.Messages)
				result.Messages = messages

				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)			
			}
		}
	}
}

func PostRestore(w http.ResponseWriter, r *http.Request) {

	config,_ := util.GetConfig(w,r)
	printConfigDebug(config)

	pluginPath := util.GetPluginPath(config.AppPlugin)
	var messages []util.Message

	if pluginPath == "" {
		var plugin string = pluginDir + "/app/" + config.AppPlugin
		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			var errMsg string = "\nApp plugin does not exist: " + plugin
	
			message := util.SetMessage("ERROR", errMsg + " " + err.Error())
			messages = append(messages, message)
	
			var result = util.SetResult(1, messages)
	
			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}
	
		var result util.Result
		result = util.ExecutePlugin(config, "app", plugin, "--action", "postRestore")
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {	
		var result util.Result
		plugin,err := util.GetAppInterface(pluginPath)
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
				result = plugin.PostRestore(config)
				messages = util.PrependMessages(setEnvResult.Messages,result.Messages)
				result.Messages = messages

				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)			
			}
		}
	}
}

func PreAppRestoreCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	config,_ := util.GetConfig(w,r)
	printConfigDebug(config)

	if config.PreAppRestoreCmd != "" {
		args := strings.Split(config.PreAppRestoreCmd, ",")
		message := util.SetMessage("INFO", "Performing pre restore app command [" + config.PreAppRestoreCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

func PostAppRestoreCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	config,_ := util.GetConfig(w,r)
	printConfigDebug(config)

	if config.PostAppRestoreCmd != "" {
		args := strings.Split(config.PostAppRestoreCmd, ",")
		message := util.SetMessage("INFO", "Performing pre restore app command [" + config.PostAppRestoreCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

