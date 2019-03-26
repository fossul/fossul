package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"engine/util"
	"net/http"
	"strings"
	"log"
)

func GetStatus(w http.ResponseWriter, r *http.Request) {
	var status = util.Status{Msg: "OK"}
	json.NewEncoder(w).Encode(status)
}

func SendTrapSuccessCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result

	var config util.Config = util.GetConfig(w,r)

	if config.SendTrapSuccessCmd != "" {
		args := strings.Split(config.SendTrapSuccessCmd, ",")
		message := util.SetMessage("INFO", "Performing send trap success command [" + config.SendTrapSuccessCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

func SendTrapErrorCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result

	var config util.Config = util.GetConfig(w,r)

	if config.SendTrapSuccessCmd != "" {
		args := strings.Split(config.SendTrapErrorCmd, ",")
		message := util.SetMessage("INFO", "Performing send trap error command [" + config.SendTrapSuccessCmd + "]")

		result = util.ExecuteCommand(args...)
		result.Messages = util.PrependMessage(message,result.Messages)

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

func GetConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var profileName string = params["profileName"]
	var configName string = params["configName"]

	conf := configDir + profileName + "/" + configName + ".conf"
	log.Println("DEBUG", "Config path is " + conf)
	var config util.Config = util.ReadConfig(conf)

	_ = json.NewDecoder(r.Body).Decode(&config)
	json.NewEncoder(w).Encode(config)
}

func GetDefaultConfig(w http.ResponseWriter, r *http.Request) {

	conf := configDir + "default" + "/" + "default.conf"
	log.Println("DEBUG", "Default config path is " + conf)
	var config util.Config = util.ReadConfig(conf)

	_ = json.NewDecoder(r.Body).Decode(&config)
	json.NewEncoder(w).Encode(config)
}

func GetDefaultPluginConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var pluginName string = params["pluginName"]

	conf := configDir + "default" + "/" + pluginName + ".conf"
	log.Println("DEBUG", "Config path is " + conf)
	configMap := util.ReadConfigToMap(conf)

	_ = json.NewDecoder(r.Body).Decode(&configMap)
	json.NewEncoder(w).Encode(configMap)
}