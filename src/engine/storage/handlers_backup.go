package main

import (
	"encoding/json"
	"engine/util"
	"net/http"
	"os"
	"fmt"
	"strings"
)

func Backup(w http.ResponseWriter, r *http.Request) {
	var config util.Config = util.GetConfig(w,r)
	pluginPath := util.GetPluginPath(config.StoragePlugin)
	var result util.Result
	var messages []util.Message

	if pluginPath == "" {
		var plugin string = config.PluginDir + "/storage/" + config.StoragePlugin

		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			var errMsg string = "Storage plugin does not exist"

			message := util.SetMessage("ERROR", errMsg + " " + err.Error())
			messages = append(messages, message)

			result = util.SetResult(1, messages)
			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}
		result = util.ExecutePlugin(config, "storage", plugin, "--action", "backup")	
		_ 	= json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		plugin := util.GetStorageInterface(pluginPath)
		setEnvResult := plugin.SetEnv(config)
		if setEnvResult.Code != 0 {
			_ = json.NewDecoder(r.Body).Decode(&setEnvResult)
			json.NewEncoder(w).Encode(setEnvResult)
		} else {
			result = plugin.Backup()
			for _,msg := range setEnvResult.Messages {
				messages = util.PrependMessage(msg,result.Messages)
			}

			if len(messages) != 0 {
				result.Messages = messages		
			}

			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)			
		}		
	}	
}

func BackupList(w http.ResponseWriter, r *http.Request) {
	var config util.Config = util.GetConfig(w,r)
	pluginPath := util.GetPluginPath(config.StoragePlugin)
	var result util.ResultSimple

	if pluginPath == "" {
		var plugin string = config.PluginDir + "/storage/" + config.StoragePlugin

		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			var errMsg string = "Storage plugin does not exist"

			var messages []string
			message := fmt.Sprintf("ERROR %s %s",errMsg,err.Error())
			messages = append(messages, message)

			var result = util.SetResultSimple(1, messages)

			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}
		result = util.ExecutePluginSimple(config, "storage", plugin, "--action", "backupList")

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		plugin := util.GetStorageInterface(pluginPath)
		_= plugin.SetEnv(config)

		backupList := plugin.BackupList()
		b, err := json.Marshal(backupList)
		if err != nil {
			result.Code = 1
			result.Messages = append(result.Messages,err.Error())
		} else {
			result.Code = 0
			outputArray := strings.Split(string(b), "\n")
			result.Messages = outputArray
		}
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)						
	}	
}

func BackupDelete(w http.ResponseWriter, r *http.Request) {
	var config util.Config = util.GetConfig(w,r)
	pluginPath := util.GetPluginPath(config.StoragePlugin)
	var result util.Result
	var messages []util.Message

	if pluginPath == "" {
		var plugin string = config.PluginDir + "/storage/" + config.StoragePlugin
		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			var errMsg string = "Storage plugin does not exist"

			var messages []util.Message
			message := util.SetMessage("ERROR", errMsg + " " + err.Error())
			messages = append(messages, message)

			result = util.SetResult(1, messages)
			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}
		result = util.ExecutePlugin(config, "storage", plugin, "--action", "backupDelete")
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		plugin := util.GetStorageInterface(pluginPath)
		setEnvResult := plugin.SetEnv(config)
		if setEnvResult.Code != 0 {
			_ = json.NewDecoder(r.Body).Decode(&setEnvResult)
			json.NewEncoder(w).Encode(setEnvResult)
		} else {
			result = plugin.BackupDelete()
			for _,msg := range setEnvResult.Messages {
				messages = util.PrependMessage(msg,result.Messages)
			}

			if len(messages) != 0 {
				result.Messages = messages		
			}

			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)			
		}			
	}	
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