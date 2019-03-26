package main

import (
	"encoding/json"
	"engine/util"
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"io/ioutil"
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
	var plugin string = config.PluginDir + "/" + pluginType + "/" + pluginName

	var result util.ResultSimple
	result = util.ExecutePluginSimple(config, pluginType, plugin, "--action", "info")

	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}