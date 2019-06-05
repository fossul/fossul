package main

import (
	"encoding/json"
	"fossul/src/engine/client/k8s"
	"fossul/src/engine/util"
	"net/http"
	"os"
	"strings"
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

			var result = util.SetResult(1, messages)

			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}

		result = util.ExecutePlugin(config, "app", plugin, "--action", "preRestore")
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		plugin, err := util.GetAppInterface(pluginPath)
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
				messages = util.PrependMessages(setEnvResult.Messages, result.Messages)
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
		}

		result = util.ExecutePlugin(config, "app", plugin, "--action", "postRestore")
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
				result = plugin.PostRestore(config)
				messages = util.PrependMessages(setEnvResult.Messages, result.Messages)
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

	if config.PreAppRestoreCmd != "" {
		args := strings.Split(config.PreAppRestoreCmd, ",")
		var messages []util.Message

		if k8s.IsRemoteCommand(args[0]) {
			args[0] = strings.Replace(args[0], ":", "", 1)
			podName, err := k8s.GetPod(config.AppPluginParameters["Namespace"], config.AppPluginParameters["ServiceName"], config.AppPluginParameters["AccessWithinCluster"])
			if err != nil {
				msg := util.SetMessage("ERROR", err.Error())
				messages = append(messages, msg)

				result = util.SetResult(1, messages)
				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)
			}

			message := util.SetMessage("INFO", "Performing remote pre restore app command ["+config.PreAppRestoreCmd+"] on pod ["+podName+"]")
			messages = append(messages, message)

			cmdResult := k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AppPluginParameters["AccessWithinCluster"], args...)

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
			message := util.SetMessage("INFO", "Performing pre restore app command ["+config.PreAppRestoreCmd+"]")

			result = util.ExecuteCommand(args...)
			result.Messages = util.PrependMessage(message, result.Messages)

			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}		
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

	if config.PostAppRestoreCmd != "" {
		args := strings.Split(config.PostAppRestoreCmd, ",")
		var messages []util.Message

		if k8s.IsRemoteCommand(args[0]) {
			args[0] = strings.Replace(args[0], ":", "", 1)
			podName, err := k8s.GetPod(config.AppPluginParameters["Namespace"], config.AppPluginParameters["ServiceName"], config.AppPluginParameters["AccessWithinCluster"])
			if err != nil {
				msg := util.SetMessage("ERROR", err.Error())
				messages = append(messages, msg)

				result = util.SetResult(1, messages)
				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)
			}

			message := util.SetMessage("INFO", "Performing remote post restore app command ["+config.PostAppRestoreCmd+"] on pod ["+podName+"]")
			messages = append(messages, message)

			cmdResult := k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AppPluginParameters["AccessWithinCluster"], args...)

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
			message := util.SetMessage("INFO", "Performing post restore app command ["+config.PostAppRestoreCmd+"]")

			result = util.ExecuteCommand(args...)
			result.Messages = util.PrependMessage(message, result.Messages)

			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}		
	}

}
