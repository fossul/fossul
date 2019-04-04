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
	pluginPath := util.GetPluginPath(config.ArchivePlugin)
	var result util.Result
	var messages []util.Message

	if pluginPath == "" {
		var plugin string = config.PluginDir + "/archive/" + config.ArchivePlugin

		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			var errMsg string = "ERROR: Archive plugin does not exist"
			log.Println(err, errMsg)

			message := util.SetMessage("ERROR", errMsg + " " + err.Error())
			messages = append(messages, message)

			var result = util.SetResult(1, messages)

			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}

		result = util.ExecutePlugin(config, "archive", plugin, "--action", "archive")	
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		plugin,err := util.GetArchiveInterface(pluginPath)
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
				result = plugin.Archive()
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
}

func ArchiveList(w http.ResponseWriter, r *http.Request) {
	var config util.Config = util.GetConfig(w,r)
	pluginPath := util.GetPluginPath(config.ArchivePlugin)
	var result util.ResultSimple
	var messages []string

	if pluginPath == "" {
		var plugin string = config.PluginDir + "/archive/" + config.ArchivePlugin

		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			var errMsg string = "ERROR: Archive plugin does not exist"
			log.Println(err, errMsg)

			message := fmt.Sprintf("ERROR %s %s",errMsg,err.Error())
			messages = append(messages, message)

			var result = util.SetResultSimple(1, messages)

			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}

		result = util.ExecutePluginSimple(config, "archive", plugin, "--action", "archiveList")

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		plugin,err := util.GetArchiveInterface(pluginPath)
		//need to implement proper result object for list calls
		if err != nil {

		} else {
			_= plugin.SetEnv(config)

			archiveList := plugin.ArchiveList()
			b, err := json.Marshal(archiveList)
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
}

func ArchiveDelete(w http.ResponseWriter, r *http.Request) {
	var config util.Config = util.GetConfig(w,r)
	pluginPath := util.GetPluginPath(config.ArchivePlugin)
	var result util.Result
	var messages []util.Message

	if pluginPath == "" {
		var plugin string = config.PluginDir + "/archive/" + config.ArchivePlugin
		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			var errMsg string = "ERROR: Archive plugin does not exist"
			log.Println(err, errMsg)

			message := util.SetMessage("ERROR", errMsg + " " + err.Error())
			messages = append(messages, message)

			result = util.SetResult(1, messages)
			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}

		result = util.ExecutePlugin(config, "archive", plugin, "--action", "archiveDelete")
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		plugin,err := util.GetArchiveInterface(pluginPath)
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
				result = plugin.ArchiveDelete()
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