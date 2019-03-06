package main

import (
	"encoding/json"
	"engine/util"
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"os"
	"io/ioutil"
)

func GetStatus(w http.ResponseWriter, r *http.Request) {
	var status = util.Status{Msg: "OK"}
	json.NewEncoder(w).Encode(status)
}

func ListPlugins(w http.ResponseWriter, r *http.Request) {
	var config util.Config = util.GetConfig(w,r)

	var plugins []string
	fileInfo, err := ioutil.ReadDir(config.PluginDir)
	if err != nil {
		log.Println(err)
	}

	for _, file := range fileInfo {
		plugins = append(plugins, file.Name())
	}

	_ = json.NewDecoder(r.Body).Decode(&plugins)
	json.NewEncoder(w).Encode(plugins)
}

func ListPluginCapabilities(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var pluginName string = params["plugin"]

	var config util.Config = util.GetConfig(w,r)
	var plugin string = config.PluginDir + "/" + pluginName

	var result util.Result
	result = util.ExecuteCommand(plugin, "--action", "list")

	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func PreQuiesceCmd(w http.ResponseWriter, r *http.Request) {

	var result = util.SetResult(0, "pre quiesce cmd completed successfully", "executed command xyz successfully")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func QuiesceCmd(w http.ResponseWriter, r *http.Request) {

	var result = util.SetResult(0, "quiesce cmd completed successfully", "executed command xyz successfully")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func Quiesce(w http.ResponseWriter, r *http.Request) {

	var config util.Config = util.GetConfig(w,r)
	var plugin string = config.PluginDir + "/" + config.AppPlugin
	if _, err := os.Stat(plugin); os.IsNotExist(err) {
		var errMsg string = "ERROR: App plugin does not exist"
		log.Println(err, errMsg)
		var result = util.SetResult(1, errMsg, err.Error())

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}

	var result util.Result
	result = util.ExecuteCommand(plugin, "--action", "quiesce")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func PostQuiesceCmd(w http.ResponseWriter, r *http.Request) {

	var result = util.SetResult(0, "post quiesce cmd completed successfully", "executed command xyz successfully")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func UnquiesceCmd(w http.ResponseWriter, r *http.Request) {

	var result = util.SetResult(0, "unquiesce cmd completed successfully", "executed command xyz successfully")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func PreUnquiesceCmd(w http.ResponseWriter, r *http.Request) {

	var result = util.SetResult(0, "pre unquiesce cmd completed successfully", "executed command xyz successfully")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func Unquiesce(w http.ResponseWriter, r *http.Request) {

	var config util.Config = util.GetConfig(w,r)
	var plugin string = config.PluginDir + "/" + config.AppPlugin
	if _, err := os.Stat(plugin); os.IsNotExist(err) {
		var errMsg string = "ERROR: App plugin does not exist"
		log.Println(err, errMsg)
		var result = util.SetResult(1, errMsg, err.Error())

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}

	var result util.Result
	result = util.ExecuteCommand(plugin, "--action", "unquiesce")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func PostUnquiesceCmd(w http.ResponseWriter, r *http.Request) {

	var result = util.SetResult(0, "post unquiesce cmd completed successfully", "executed command xyz successfully")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}
