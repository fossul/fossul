package main

import (
	"encoding/json"
	"fossul/src/engine/util"
	"net/http"
	"os"
	"strings"
)

// UnquiesceCmd godoc
// @Description Application unquiesce command
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /unquiesceCmd [post]
func UnquiesceCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var messages []util.Message

	config, err := util.GetConfig(w, r)
	printConfigDebug(config)

	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't read config! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	if config.AppUnquiesceCmd != "" {
		args := strings.Split(config.AppUnquiesceCmd, ",")
		message := util.SetMessage("INFO", "Performing unquiesce command ["+config.AppUnquiesceCmd+"]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message, result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// PreUnquiesceCmd godoc
// @Description Application pre unquiesce command
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /preUnquiesceCmd [post]
func PreUnquiesceCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var messages []util.Message

	config, err := util.GetConfig(w, r)
	printConfigDebug(config)

	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't read config! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	if config.PreAppUnquiesceCmd != "" {
		args := strings.Split(config.PreAppUnquiesceCmd, ",")
		message := util.SetMessage("INFO", "Performing pre unquiesce command ["+config.PreAppUnquiesceCmd+"]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message, result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// Unquiesce godoc
// @Description Application unquiesce
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /unquiesce [post]
func Unquiesce(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var messages []util.Message

	config, err := util.GetConfig(w, r)
	printConfigDebug(config)

	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't read config! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	pluginPath := util.GetPluginPath(config.AppPlugin)

	if pluginPath == "" {
		var plugin string = pluginDir + "/app/" + config.AppPlugin
		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			var errMsg string = "\nApp plugin does not exist: " + plugin

			message := util.SetMessage("ERROR", errMsg+" "+err.Error())
			messages = append(messages, message)

			result = util.SetResult(1, messages)

			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)

			return
		}

		result = util.ExecutePlugin(config, "app", plugin, "--action", "unquiesce")
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		plugin, err := util.GetAppInterface(pluginPath)
		if err != nil {
			message := util.SetMessage("ERROR", err.Error())
			messages = append(messages, message)

			result = util.SetResult(1, messages)
			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		} else {
			setEnvResult := plugin.SetEnv(config)
			if setEnvResult.Code != 0 {
				_ = json.NewDecoder(r.Body).Decode(&setEnvResult)
				json.NewEncoder(w).Encode(setEnvResult)
			} else {
				result = plugin.Unquiesce(config)
				messages = util.PrependMessages(setEnvResult.Messages, result.Messages)
				result.Messages = messages

				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)
			}
		}
	}
}

// PostUnquiesceCmd godoc
// @Description Application post unquiesce command
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /postUnquiesceCmd [post]
func PostUnquiesceCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var messages []util.Message

	config, err := util.GetConfig(w, r)
	printConfigDebug(config)

	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't read config! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	if config.PostAppUnquiesceCmd != "" {
		args := strings.Split(config.PostAppUnquiesceCmd, ",")
		message := util.SetMessage("INFO", "Performing post unquiesce command ["+config.PostAppUnquiesceCmd+"]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message, result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}
