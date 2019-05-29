package main

import (
	"encoding/json"
	"fossil/src/engine/util"
	"net/http"
	"strings"
	"os"
)

// PreRestore godoc
// @Description Application pre restore
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /preRestore [post]
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

// PostRestpore godoc
// @Description Application post restore
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /postRestore [post]
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

// PreAppRestoreCmd godoc
// @Description Application pre restore command
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /preAppRestoreCmd [post]
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

// PostAppRestoreCmd godoc
// @Description Application post restore command
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /postAppRestoreCmd [post]
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

