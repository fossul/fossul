package main

import (
	"encoding/json"
	"engine/util"
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"io/ioutil"
	"strings"
	"os"
)

func GetStatus(w http.ResponseWriter, r *http.Request) {
	var status = util.Status{Msg: "OK"}
	json.NewEncoder(w).Encode(status)
}

func PluginList(w http.ResponseWriter, r *http.Request) {
	var plugins []string

	params := mux.Vars(r)	
	var pluginType string = params["pluginType"]

	var config util.Config = util.GetConfig(w,r)
	var pluginDir string

	if pluginType == "app" {
		pluginDir = config.PluginDir + "/app"
	} else {
		log.Println("ERROR plugin type " + pluginType + " must be app")

		_ = json.NewDecoder(r.Body).Decode(&plugins)
		json.NewEncoder(w).Encode(plugins)
	}

	fileInfo, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		log.Println(err)
	}

	for _, file := range fileInfo {
		fileName := file.Name()
		plugins = append(plugins, fileName)
	}

	_ = json.NewDecoder(r.Body).Decode(&plugins)
	json.NewEncoder(w).Encode(plugins)
}

func PluginInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var pluginName string = params["pluginName"]
	var pluginType string = params["pluginType"]

	var pluginInfoResult util.PluginInfoResult
	var pluginInfo util.Plugin
	var result util.Result
	var messages []util.Message

	var config util.Config = util.GetConfig(w,r)
	pluginPath := util.GetPluginPath(pluginName)

	if pluginPath == "" {
		var plugin string = config.PluginDir + "/" + pluginType + "/" + pluginName
		if _, err := os.Stat(plugin); os.IsNotExist(err) {
			msg := util.SetMessage("ERROR","Plugin not found! " + err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			pluginInfoResult.Result = result

			_ = json.NewDecoder(r.Body).Decode(&pluginInfoResult)
			json.NewEncoder(w).Encode(pluginInfoResult)
		}

		var resultSimple util.ResultSimple
		resultSimple = util.ExecutePluginSimple(config, pluginType, plugin, "--action", "info")
		if resultSimple.Code != 0 {
			msg := util.SetMessage("ERROR","Plugin Info failed!")
			messages = append(messages, msg)	
			result := util.SetResult(1, messages)
			pluginInfoResult.Result = result

			_ = json.NewDecoder(r.Body).Decode(&pluginInfoResult)
			json.NewEncoder(w).Encode(pluginInfoResult)
		} else {
			pluginInfoString := strings.Join(resultSimple.Messages," ")
			json.Unmarshal([]byte(pluginInfoString), &pluginInfo)
		
			pluginInfoResult.Result.Code = resultSimple.Code
			pluginInfoResult.Plugin = pluginInfo
	
			_ = json.NewDecoder(r.Body).Decode(&pluginInfoResult)
			json.NewEncoder(w).Encode(pluginInfoResult)
		}
	} else {
		plugin,err := util.GetAppInterface(pluginPath)

		if err != nil {	
			msg := util.SetMessage("ERROR",err.Error())
			messages = append(messages,msg)
			result = util.SetResult(1,messages)
			pluginInfoResult.Result = result	
				
			_ = json.NewDecoder(r.Body).Decode(&pluginInfoResult)
			json.NewEncoder(w).Encode(pluginInfoResult)	
		} else {
			setEnvResult := plugin.SetEnv(config)
			if setEnvResult.Code != 0 {
				pluginInfoResult.Result = setEnvResult
				_ = json.NewDecoder(r.Body).Decode(&pluginInfoResult)
				json.NewEncoder(w).Encode(pluginInfoResult)
			} else {	
				pluginInfo := plugin.Info()
				
				pluginInfoResult.Result.Code = 0
				pluginInfoResult.Plugin = pluginInfo

				_ = json.NewDecoder(r.Body).Decode(&pluginInfoResult)
				json.NewEncoder(w).Encode(pluginInfoResult)		
			}	
		}			
	}	
}