package main

import (
	"encoding/json"
	"fossil/src/engine/util"
	"net/http"
	"os"
	"strings"
)

// PreQuiesceCmd godoc
// @Description Application pre quiesce command
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /preQuiesceCmd [post]
func PreQuiesceCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	config,_ := util.GetConfig(w,r)
	printConfigDebug(config)

	if config.PreAppQuiesceCmd != "" {
		args := strings.Split(config.PreAppQuiesceCmd, ",")
		message := util.SetMessage("INFO", "Performing pre quiesce command [" + config.PreAppQuiesceCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// QuiesceCmd godoc
// @Description Application quiesce command
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /puiesceCmd [post]
func QuiesceCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	config,_ := util.GetConfig(w,r)
	printConfigDebug(config)

	if config.PreAppQuiesceCmd != "" {
		args := strings.Split(config.AppQuiesceCmd, ",")
		message := util.SetMessage("INFO", "Performing quiesce command [" + config.PreAppQuiesceCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// Quiesce godoc
// @Description Application quiesce
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /quiesce [post]
func Quiesce(w http.ResponseWriter, r *http.Request) {

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
		result = util.ExecutePlugin(config, "app", plugin, "--action", "quiesce")
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
				result = plugin.Quiesce(config)
				messages = util.PrependMessages(setEnvResult.Messages,result.Messages)
				result.Messages = messages

				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)			
			}
		}
	}
}

// PostQuiesceCmd godoc
// @Description Application post quiesce command
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /postQuiesceCmd [post]
func PostQuiesceCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	config,_ := util.GetConfig(w,r)
	printConfigDebug(config)

	if config.PreAppQuiesceCmd != "" {
		args := strings.Split(config.PostAppQuiesceCmd, ",")
		message := util.SetMessage("INFO", "Performing post quiesce command [" + config.PreAppQuiesceCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}