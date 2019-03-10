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

func PluginList(w http.ResponseWriter, r *http.Request) {
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

func PluginInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var pluginName string = params["plugin"]

	var config util.Config = util.GetConfig(w,r)
	var plugin string = config.PluginDir + "/" + pluginName

	var result util.Result
	result = util.ExecutePlugin(config, "app", plugin, "--action", "info")

	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func PreQuiesceCmd(w http.ResponseWriter, r *http.Request) {

	messages := []string{"pre quiesce cmd completed successfully","executed command xyz successfully"}
	var result = util.SetResult(0, messages)
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func QuiesceCmd(w http.ResponseWriter, r *http.Request) {

	messages := []string{"quiesce cmd completed successfully","executed command xyz successfully"}
	var result = util.SetResult(0, messages)
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func Quiesce(w http.ResponseWriter, r *http.Request) {

	var config util.Config = util.GetConfig(w,r)
	var plugin string = config.PluginDir + "/" + config.AppPlugin
	if _, err := os.Stat(plugin); os.IsNotExist(err) {
		var errMsg string = "ERROR: App plugin does not exist"
		log.Println(err, errMsg)
		messages := []string{errMsg,err.Error()}
		var result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}

	var result util.Result
	result = util.ExecutePlugin(config, "app", plugin, "--action", "quiesce")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func PostQuiesceCmd(w http.ResponseWriter, r *http.Request) {

	messages := []string{"post quiesce cmd completed successfully","executed command xyz successfully"}
	var result = util.SetResult(0, messages)
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func UnquiesceCmd(w http.ResponseWriter, r *http.Request) {

	messages := []string{"unquiesce cmd completed successfully","executed command xyz successfully"}
	var result = util.SetResult(0, messages)
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func PreUnquiesceCmd(w http.ResponseWriter, r *http.Request) {

	messages := []string{"pre unquiesce cmd completed successfully","executed command xyz successfully"}
	var result = util.SetResult(0, messages)
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func Unquiesce(w http.ResponseWriter, r *http.Request) {

	var config util.Config = util.GetConfig(w,r)
	var plugin string = config.PluginDir + "/" + config.AppPlugin
	
	if _, err := os.Stat(plugin); os.IsNotExist(err) {
		var errMsg string = "ERROR: App plugin does not exist"
		log.Println(err, errMsg)
		messages := []string{errMsg,err.Error()}
		var result = util.SetResult(1, messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}

	var result util.Result
	result = util.ExecutePlugin(config, "app", plugin, "--action", "unquiesce")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func PostUnquiesceCmd(w http.ResponseWriter, r *http.Request) {

	messages := []string{"post unquiesce cmd completed successfully","executed command xyz successfully"}
	var result = util.SetResult(0, messages)
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}
