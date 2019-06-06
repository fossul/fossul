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

// Discover godoc
// @Description Application discover
// @Param config body util.Config true "config struct"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /discover [post]
func Discover(w http.ResponseWriter, r *http.Request) {
	var discoverResult util.DiscoverResult
	var result util.Result
	var messages []util.Message

	config, err := util.GetConfig(w, r)
	printConfigDebug(config)

	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't read config! "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)

		discoverResult.Result = result

		_ = json.NewDecoder(r.Body).Decode(&discoverResult)
		json.NewEncoder(w).Encode(discoverResult)

		return
	}

	pluginPath := util.GetPluginPath(config.AppPlugin)
	//
	if pluginPath == "" {
		var plugin string = pluginDir + "/app/" + config.AppPlugin
		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			var errMsg string = "\nApp plugin does not exist: " + plugin

			message := util.SetMessage("ERROR", errMsg+" "+err.Error())
			messages = append(messages, message)

			result = util.SetResult(1, messages)

			discoverResult.Result = result
			_ = json.NewDecoder(r.Body).Decode(&discoverResult)
			json.NewEncoder(w).Encode(discoverResult)
		}
		resultSimple := util.ExecutePluginSimple(config, "app", plugin, "--action", "discover")
		discoverResultString := strings.Join(resultSimple.Messages, " ")
		json.Unmarshal([]byte(discoverResultString), &discoverResult)

		_ = json.NewDecoder(r.Body).Decode(&discoverResult)
		json.NewEncoder(w).Encode(discoverResult)
	} else {
		plugin, err := util.GetAppInterface(pluginPath)
		if err != nil {
			message := util.SetMessage("ERROR", err.Error())
			messages = append(messages, message)

			result = util.SetResult(1, messages)
			discoverResult.Result = result
			_ = json.NewDecoder(r.Body).Decode(&discoverResult)
			json.NewEncoder(w).Encode(discoverResult)
		} else {
			setEnvResult := plugin.SetEnv(config)
			if setEnvResult.Code != 0 {
				discoverResult.Result = setEnvResult
				_ = json.NewDecoder(r.Body).Decode(&discoverResult)
				json.NewEncoder(w).Encode(discoverResult)
			} else {
				discoverResult = plugin.Discover(config)
				messages = util.PrependMessages(setEnvResult.Messages, discoverResult.Result.Messages)
				discoverResult.Result.Messages = messages

				_ = json.NewDecoder(r.Body).Decode(&discoverResult)
				json.NewEncoder(w).Encode(discoverResult)
			}
		}
	}
}
