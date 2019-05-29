package main

import (
	"encoding/json"
	"fossil/src/engine/util"
	"net/http"
	"os"
	"strings"
)

// Restore godoc
// @Description Restore data
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /restore [post]
func Restore(w http.ResponseWriter, r *http.Request) {
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

// RestoreCmd godoc
// @Description Restore Command
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /restoreCmd [post]
func RestoreCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	config,_ := util.GetConfig(w,r)
	printConfigDebug(config)

	if config.RestoreCmd != "" {
		args := strings.Split(config.RestoreCmd, ",")
		message := util.SetMessage("INFO", "Performing restore command [" + config.RestoreCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}