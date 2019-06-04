package main

import (
	"encoding/json"
	"fmt"
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

		result = util.ExecutePlugin(config, "archive", plugin, "--action", "archive")
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
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /archiveList [post]
func ArchiveList(w http.ResponseWriter, r *http.Request) {
	config, _ := util.GetConfig(w, r)
	printConfigDebug(config)

	pluginPath := util.GetPluginPath(config.ArchivePlugin)
	var result util.ResultSimple
	var messages []string

	if pluginPath == "" {
		var plugin string = pluginDir + "/archive/" + config.ArchivePlugin

		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			var errMsg string = "Archive plugin does not exist"

			message := fmt.Sprintf("ERROR %s %s", errMsg, err.Error())
			messages = append(messages, message)

			var result = util.SetResultSimple(1, messages)

			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)
		}

		result = util.ExecutePluginSimple(config, "archive", plugin, "--action", "archiveList")

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		plugin, err := util.GetArchiveInterface(pluginPath)
		//need to implement proper result object for list calls
		if err != nil {

		} else {
			_ = plugin.SetEnv(config)

			archiveList := plugin.ArchiveList(config)
			b, err := json.Marshal(archiveList)
			if err != nil {
				result.Code = 1
				result.Messages = append(result.Messages, err.Error())
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

		result = util.ExecutePlugin(config, "archive", plugin, "--action", "archiveDelete")
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

	if config.BackupCreateCmd != "" {
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

	if config.BackupDeleteCmd != "" {
		args := strings.Split(config.ArchiveDeleteCmd, ",")
		message := util.SetMessage("INFO", "Performing archive delete command ["+config.ArchiveDeleteCmd+"]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message, result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}
