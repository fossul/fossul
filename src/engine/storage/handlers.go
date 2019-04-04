package main

import (
	"encoding/json"
	"engine/util"
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"io/ioutil"
	"strings"
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

	if pluginType == "storage" {
		pluginDir = config.PluginDir + "/storage"
	} else if pluginType == "archive" {
		pluginDir = config.PluginDir + "/archive"
	} else {
		log.Println("ERROR plugin type " + pluginType + " must be storage|archive")

		_ = json.NewDecoder(r.Body).Decode(&plugins)
		json.NewEncoder(w).Encode(plugins)
	}

	fileInfo, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		log.Println(err)
	}

	for _, file := range fileInfo {
		plugins = append(plugins, file.Name())
	}

	_ = json.NewDecoder(r.Body).Decode(&plugins)
	json.NewEncoder(w).Encode(plugins)
}

func PluginInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var pluginName string = params["pluginName"]
	var pluginType string = params["pluginType"]

	var config util.Config = util.GetConfig(w,r)
	pluginPath := util.GetPluginPath(pluginName)

	if pluginPath == "" {
		var plugin string = config.PluginDir + "/" + pluginType + "/" + pluginName

		var result util.ResultSimple
		result = util.ExecutePluginSimple(config, pluginType, plugin, "--action", "info")
	
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		var result util.ResultSimple
		if pluginType == "storage" {
			plugin,err := util.GetStorageInterface(pluginPath)
			//need to implement proper result object for info
			if err != nil {	
			} else {
				plugin.SetEnv(config)
				pluginInfo := plugin.Info()
				b, err := json.Marshal(pluginInfo)
				if err != nil {
					result.Code = 1
					result.Messages = append(result.Messages,err.Error())
				} else {
					result.Code = 0
					outputArray := strings.Split(string(b), "\n")
					result.Messages = outputArray
				}	
				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)		
			}			
		} else if pluginType == "archive" {
			plugin,err := util.GetArchiveInterface(pluginPath)
			//need to implement proper result object for info
			if err != nil {			
			} else {
				plugin.SetEnv(config)
				pluginInfo := plugin.Info()
				b, err := json.Marshal(pluginInfo)
				if err != nil {
					result.Code = 1
					result.Messages = append(result.Messages,err.Error())
				} else {
					result.Code = 0
					outputArray := strings.Split(string(b), "\n")
					result.Messages = outputArray
				}		
				_ = json.NewDecoder(r.Body).Decode(&result)
				json.NewEncoder(w).Encode(result)	
			}	
		} else {
			_ = json.NewDecoder(r.Body).Decode(&result)
			json.NewEncoder(w).Encode(result)	
		}
	}	
}