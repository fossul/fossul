package main

import (
	"encoding/json"
	"fossil/src/engine/util"
	"net/http"
	"log"
	"os"
	"strings"
)

func Discover(w http.ResponseWriter, r *http.Request) {

	config,_ := util.GetConfig(w,r)
	pluginPath := util.GetPluginPath(config.AppPlugin)
	var discoverResult util.DiscoverResult
	var messages []util.Message

	if pluginPath == "" {
		var plugin string = pluginDir + "/app/" + config.AppPlugin
		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			var errMsg string = "\nERROR: App plugin does not exist: " + plugin
			log.Println(err, errMsg)
	
			message := util.SetMessage("ERROR", errMsg + " " + err.Error())
			messages = append(messages, message)
	
			result := util.SetResult(1, messages)
	
			discoverResult.Result = result
			_ = json.NewDecoder(r.Body).Decode(&discoverResult)
			json.NewEncoder(w).Encode(discoverResult)
		}
		resultSimple := util.ExecutePluginSimple(config, "app", plugin, "--action", "discover")
		discoverResultString := strings.Join(resultSimple.Messages," ")
		json.Unmarshal([]byte(discoverResultString), &discoverResult)

		_ = json.NewDecoder(r.Body).Decode(&discoverResult)
		json.NewEncoder(w).Encode(discoverResult)	
	} else {	
		plugin,err := util.GetAppInterface(pluginPath)
		if err != nil {
			message := util.SetMessage("ERROR", err.Error())
			messages = append(messages, message)

			var result = util.SetResult(1, messages)
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
				discoverResult = plugin.Discover()
				messages = util.PrependMessages(setEnvResult.Messages,discoverResult.Result.Messages)
				discoverResult.Result.Messages = messages

				_ = json.NewDecoder(r.Body).Decode(&discoverResult)
				json.NewEncoder(w).Encode(discoverResult)			
			}
		}
	}
}