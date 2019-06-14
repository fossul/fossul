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
	"fossul/src/engine/util"
	"net/http"
	"os"
	"strings"
)

// Archive godoc
// @Description Archive backup
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /archive [post]
func Archive(w http.ResponseWriter, r *http.Request) {
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

	pluginPath := util.GetPluginPath(config.ArchivePlugin)

	if pluginPath == "" {
		var plugin string = pluginDir + "/archive/" + config.ArchivePlugin

		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			var errMsg string = "Archive plugin does not exist"

			message := util.SetMessage("ERROR", errMsg+" "+err.Error())
			messages = append(messages, message)

			result = util.SetResult(1, messages)

			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}

		result = util.ExecutePlugin(config, "archive", plugin, "--archive")
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		plugin, err := util.GetArchiveInterface(pluginPath)
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
				result = plugin.Archive(config)
				messages = util.PrependMessages(setEnvResult.Messages, result.Messages)
				result.Messages = messages

				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)
			}
		}
	}
}

// ArchiveList godoc
// @Description List archive backups
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Archives
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /archiveList [post]
func ArchiveList(w http.ResponseWriter, r *http.Request) {
	var archives util.Archives
	var archiveList []util.Archive
	var result util.Result
	var messages []util.Message

	config, err := util.GetConfig(w, r)
	printConfigDebug(config)

	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't read config! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		archives.Result = result

		_ = json.NewDecoder(r.Body).Decode(&archives)
		json.NewEncoder(w).Encode(archives)

		return
	}

	pluginPath := util.GetPluginPath(config.ArchivePlugin)

	if pluginPath == "" {
		var plugin string = pluginDir + "/archive/" + config.ArchivePlugin

		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			msg := util.SetMessage("ERROR", "Archive plugin not found! "+err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			archives.Result = result

			_ = json.NewDecoder(r.Body).Decode(&archives)
			json.NewEncoder(w).Encode(archives)
		}

		resultSimple := util.ExecutePluginSimple(config, "archive", plugin, "--archiveList")
		if resultSimple.Code != 0 {
			msg := util.SetMessage("ERROR", "ArchiveList failed")
			messages = append(messages, msg)
			result := util.SetResult(1, messages)
			archives.Result = result

			_ = json.NewDecoder(r.Body).Decode(&archives)
			json.NewEncoder(w).Encode(archives)
		} else {
			archiveListString := strings.Join(resultSimple.Messages, " ")
			json.Unmarshal([]byte(archiveListString), &archiveList)

			archives.Result.Code = resultSimple.Code
			archives.Archives = archiveList

			_ = json.NewDecoder(r.Body).Decode(&archives)
			json.NewEncoder(w).Encode(archives)
		}
	} else {
		plugin, err := util.GetArchiveInterface(pluginPath)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)
			result = util.SetResult(1, messages)
			archives.Result = result

			_ = json.NewDecoder(r.Body).Decode(&archives)
			json.NewEncoder(w).Encode(archives)
		} else {
			setEnvResult := plugin.SetEnv(config)
			if setEnvResult.Code != 0 {
				archives.Result = setEnvResult
				_ = json.NewDecoder(r.Body).Decode(&archives)
				json.NewEncoder(w).Encode(archives)
			} else {
				archives := plugin.ArchiveList(config)
				_ = json.NewDecoder(r.Body).Decode(&archives)
				json.NewEncoder(w).Encode(archives)
			}
		}
	}
}

// ArchiveDelete godoc
// @Description Delete archive backups according to retention
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /archiveDelete [post]
func ArchiveDelete(w http.ResponseWriter, r *http.Request) {
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

	pluginPath := util.GetPluginPath(config.ArchivePlugin)

	if pluginPath == "" {
		var plugin string = pluginDir + "/archive/" + config.ArchivePlugin
		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			var errMsg string = "Archive plugin does not exist"

			message := util.SetMessage("ERROR", errMsg+" "+err.Error())
			messages = append(messages, message)

			result = util.SetResult(1, messages)
			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}

		result = util.ExecutePlugin(config, "archive", plugin, "--archiveDelete")
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		plugin, err := util.GetArchiveInterface(pluginPath)
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
				result = plugin.ArchiveDelete(config)
				messages = util.PrependMessages(setEnvResult.Messages, result.Messages)
				result.Messages = messages

				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)
			}
		}
	}
}

// ArchiveCreateCmd godoc
// @Description Create archive command
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /archiveCreateCmd [post]
func ArchiveCreateCmd(w http.ResponseWriter, r *http.Request) {
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

	if config.ArchiveCreateCmd != "" {
		args := strings.Split(config.ArchiveCreateCmd, ",")
		message := util.SetMessage("INFO", "Performing archive create command ["+config.ArchiveCreateCmd+"]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message, result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// ArchiveDeleteCmd godoc
// @Description Delete archive command
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /archiveDeleteCmd [post]
func ArchiveDeleteCmd(w http.ResponseWriter, r *http.Request) {
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

	if config.ArchiveDeleteCmd != "" {
		args := strings.Split(config.ArchiveDeleteCmd, ",")
		message := util.SetMessage("INFO", "Performing archive delete command ["+config.ArchiveDeleteCmd+"]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message, result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}
