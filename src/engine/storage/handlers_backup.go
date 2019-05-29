package main

import (
	"encoding/json"
	"fossil/src/engine/util"
	"net/http"
	"os"
	"strings"
)

// Backup godoc
// @Description Backup data
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /backup [post]
func Backup(w http.ResponseWriter, r *http.Request) {
	config,_ := util.GetConfig(w,r)
	printConfigDebug(config)

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
		result = util.ExecutePlugin(config, "storage", plugin, "--action", "backup")	
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
				result = plugin.Backup(config)
				messages = util.PrependMessages(setEnvResult.Messages,result.Messages)
				result.Messages = messages
	
				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)			
			}	
		}	
	}	
}

// BackupList godoc
// @Description List backups
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /backupList [post]
func BackupList(w http.ResponseWriter, r *http.Request) {
	config,_ := util.GetConfig(w,r)
	printConfigDebug(config)

	pluginPath := util.GetPluginPath(config.StoragePlugin)

	var backups util.Backups
	var result util.Result
	var messages []util.Message
	if pluginPath == "" {
		var plugin string = pluginDir + "/storage/" + config.StoragePlugin

		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			msg := util.SetMessage("ERROR","Storage plugin not found! " + err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			backups.Result = result

			_ = json.NewDecoder(r.Body).Decode(&backups)
			json.NewEncoder(w).Encode(backups)
		}
		var backups util.Backups
		var backupList []util.Backup
		resultSimple := util.ExecutePluginSimple(config, "storage", plugin, "--action", "backupList")
		if resultSimple.Code != 0 {
			msg := util.SetMessage("ERROR","BackupList failed")
			messages = append(messages, msg)	
			result := util.SetResult(1, messages)
			backups.Result = result

			_ = json.NewDecoder(r.Body).Decode(&backups)
			json.NewEncoder(w).Encode(backups)
		} else {
			backupListString := strings.Join(resultSimple.Messages," ")
			json.Unmarshal([]byte(backupListString), &backupList)
		
			backups.Result.Code = resultSimple.Code
			backups.Backups = backupList
	
			_ = json.NewDecoder(r.Body).Decode(&backups)
			json.NewEncoder(w).Encode(backups)
		}
	} else {
		plugin,err := util.GetStorageInterface(pluginPath)
		if err != nil {
			msg := util.SetMessage("ERROR",err.Error())
			messages = append(messages,msg)
			result = util.SetResult(1,messages)
			backups.Result = result

			_ = json.NewDecoder(r.Body).Decode(&backups)
			json.NewEncoder(w).Encode(backups)	
		} else {
			setEnvResult := plugin.SetEnv(config)
			if setEnvResult.Code != 0 {
				backups.Result = setEnvResult
				_ = json.NewDecoder(r.Body).Decode(&backups)
				json.NewEncoder(w).Encode(backups)
			} else {	
				backups := plugin.BackupList(config)
				_ = json.NewDecoder(r.Body).Decode(&backups)
				json.NewEncoder(w).Encode(backups)	
			}		
		}				
	}	
}

// BackupDelete godoc
// @Description Delete backups according to retention
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /backupDelete [post]
func BackupDelete(w http.ResponseWriter, r *http.Request) {
	config,_ := util.GetConfig(w,r)
	printConfigDebug(config)

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
		result = util.ExecutePlugin(config, "storage", plugin, "--action", "backupDelete")
		_ = json.NewDecoder(r.Body).Decode(&result)
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
				result = plugin.BackupDelete(config)
				messages = util.PrependMessages(setEnvResult.Messages,result.Messages)
				result.Messages = messages
	
				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)			
			}	
		}		
	}	
}

// BackupDeleteCmd godoc
// @Description Delete backup command
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /backupDeleteCmd [post]
func BackupCreateCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	config,_ := util.GetConfig(w,r)
	printConfigDebug(config)

	if config.BackupCreateCmd != "" {
		args := strings.Split(config.BackupCreateCmd, ",")
		message := util.SetMessage("INFO", "Performing backup create command [" + config.BackupCreateCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// BackupCreateCmd godoc
// @Description Create backup command
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /backupCreateCmd [post]
func BackupDeleteCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	config,_ := util.GetConfig(w,r)
	printConfigDebug(config)

	if config.BackupDeleteCmd != "" {
		args := strings.Split(config.BackupDeleteCmd, ",")
		message := util.SetMessage("INFO", "Performing backup delete command [" + config.BackupDeleteCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}