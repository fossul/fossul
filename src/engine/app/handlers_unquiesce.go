/*
Copyright 2019 The Fossul Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/fossul/fossul/src/client/k8s"
	"github.com/fossul/fossul/src/engine/util"
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
		var messages []util.Message

		if k8s.IsRemoteCommand(args[0]) {
			args[0] = strings.Replace(args[0], ":", "", 1)
			podName, err := k8s.GetPodName(config.AppPluginParameters["Namespace"], config.AppPluginParameters["PodSelector"], config.AccessWithinCluster)
			if err != nil {
				msg := util.SetMessage("ERROR", err.Error())
				messages = append(messages, msg)

				result = util.SetResult(1, messages)
				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)
			}

			message := util.SetMessage("INFO", "Performing remote unquiesce command ["+config.AppUnquiesceCmd+"] on pod ["+podName+"]")
			messages = append(messages, message)

			cmdResult := k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AccessWithinCluster, args...)

			if cmdResult.Code != 0 {
				messages = util.PrependMessages(messages, cmdResult.Messages)
				result = util.SetResult(1, messages)

				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)
			} else {
				messages = util.PrependMessages(messages, cmdResult.Messages)
				result = util.SetResult(0, messages)

				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)
			}

		} else {
			message := util.SetMessage("INFO", "Performing unquiesce command ["+config.AppUnquiesceCmd+"]")

			result = util.ExecuteCommand(args...)
			result.Messages = util.PrependMessage(message, result.Messages)

			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}
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
		var messages []util.Message

		if k8s.IsRemoteCommand(args[0]) {
			args[0] = strings.Replace(args[0], ":", "", 1)
			podName, err := k8s.GetPodName(config.AppPluginParameters["Namespace"], config.AppPluginParameters["PodSelector"], config.AccessWithinCluster)
			if err != nil {
				msg := util.SetMessage("ERROR", err.Error())
				messages = append(messages, msg)

				result = util.SetResult(1, messages)
				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)
			}

			message := util.SetMessage("INFO", "Performing remote pre unquiesce command ["+config.PreAppRestoreCmd+"] on pod ["+podName+"]")
			messages = append(messages, message)

			cmdResult := k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AccessWithinCluster, args...)

			if cmdResult.Code != 0 {
				messages = util.PrependMessages(messages, cmdResult.Messages)
				result = util.SetResult(1, messages)

				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)
			} else {
				messages = util.PrependMessages(messages, cmdResult.Messages)
				result = util.SetResult(0, messages)

				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)
			}

		} else {
			message := util.SetMessage("INFO", "Performing pre unquiesce command ["+config.PreAppUnquiesceCmd+"]")

			result = util.ExecuteCommand(args...)
			result.Messages = util.PrependMessage(message, result.Messages)

			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}
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

	pluginPath := util.GetPluginPath(config.AppPlugin, "app")

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

		result = util.ExecutePlugin(config, "app", plugin, "--unquiesce")
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
		var messages []util.Message

		if k8s.IsRemoteCommand(args[0]) {
			args[0] = strings.Replace(args[0], ":", "", 1)
			podName, err := k8s.GetPodName(config.AppPluginParameters["Namespace"], config.AppPluginParameters["PodSelector"], config.AccessWithinCluster)
			if err != nil {
				msg := util.SetMessage("ERROR", err.Error())
				messages = append(messages, msg)

				result = util.SetResult(1, messages)
				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)
			}

			message := util.SetMessage("INFO", "Performing remote post unquiesce command ["+config.PostAppUnquiesceCmd+"] on pod ["+podName+"]")
			messages = append(messages, message)

			cmdResult := k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AccessWithinCluster, args...)

			if cmdResult.Code != 0 {
				messages = util.PrependMessages(messages, cmdResult.Messages)
				result = util.SetResult(1, messages)

				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)
			} else {
				messages = util.PrependMessages(messages, cmdResult.Messages)
				result = util.SetResult(0, messages)

				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)
			}

		} else {
			message := util.SetMessage("INFO", "Performing post unquiesce command ["+config.PostAppUnquiesceCmd+"]")

			result = util.ExecuteCommand(args...)
			result.Messages = util.PrependMessage(message, result.Messages)

			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}
	}
}
