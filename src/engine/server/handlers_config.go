package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"fossil/src/engine/util"
	"net/http"
	"log"
	"os"
)

func GetConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var profileName string = params["profileName"]
	var configName string = params["configName"]

	conf := configDir + profileName + "/" + configName + ".conf"
	log.Println("DEBUG", "Config path is " + conf)
	config,_ := util.ReadConfig(conf)

	_ = json.NewDecoder(r.Body).Decode(&config)
	json.NewEncoder(w).Encode(config)
}

func GetPluginConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]	
	var pluginName string = params["pluginName"]

	conf := configDir + profileName + "/" + pluginName + ".conf"
	log.Println("DEBUG", "Plugin config path is " + conf)
	configMap,_ := util.ReadConfigToMap(conf)

	_ = json.NewDecoder(r.Body).Decode(&configMap)
	json.NewEncoder(w).Encode(configMap)
}

func GetDefaultConfig(w http.ResponseWriter, r *http.Request) {

	conf := configDir + "default" + "/" + "default.conf"
	log.Println("DEBUG", "Default config path is " + conf)
	config,_ := util.ReadConfig(conf)

	_ = json.NewDecoder(r.Body).Decode(&config)
	json.NewEncoder(w).Encode(config)
}

func GetDefaultPluginConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var pluginName string = params["pluginName"]

	conf := configDir + "default" + "/" + pluginName + ".conf"
	log.Println("DEBUG", "Config path is " + conf)
	configMap,_ := util.ReadConfigToMap(conf)

	_ = json.NewDecoder(r.Body).Decode(&configMap)
	json.NewEncoder(w).Encode(configMap)
}

func AddConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var profileName string = params["profileName"]
	var configName string = params["configName"]

	var result util.Result
	var messages []util.Message

	config,err := util.GetConfig(w,r)
	if err != nil {
		msg := util.SetMessage("ERROR","Couldn't get configuration. " + err.Error()) 
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages	

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)	

		return
	}

	dir := configDir + "/" + profileName
	conf := dir + "/" + configName + ".conf"

	if util.ExistsPath(dir) != true {
		msg := util.SetMessage("ERROR","Add configuration [" + conf + "] failed! Profile [" + dir + "] does not exist.") 
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages	

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		//err := util.WriteGob(conf,config)
		err := util.WriteConfig(conf,config)

		if err != nil {
			msg := util.SetMessage("ERROR","Add configuration [" + conf + "] failed!" + err.Error())
			messages = append(messages, msg)
	
			result.Code = 1
		} else {
			msg := util.SetMessage("INFO","Configuration [" + conf + "] create completed successfully")
			messages = append(messages, msg)
	
			result.Code = 0
		}
		result.Messages = messages
	
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

func AddPluginConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var profileName string = params["profileName"]
	var pluginName string = params["pluginName"]

	var result util.Result
	var messages []util.Message

	configMap,err := util.GetPluginConfig(w,r)
	if err != nil {
		msg := util.SetMessage("ERROR","Couldn't get configuration. " + err.Error()) 
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages	

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)	

		return
	}

	dir := configDir + "/" + profileName
	conf := dir + "/" + pluginName + ".conf"

	if util.ExistsPath(dir) != true {
		msg := util.SetMessage("ERROR","Add configuration [" + conf + "] failed! Profile [" + dir + "] does not exist.") 
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages	

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		//err := util.WriteGob(conf,config)
		err := util.WritePluginConfig(conf,configMap)

		if err != nil {
			msg := util.SetMessage("ERROR","Add configuration [" + conf + "] failed!" + err.Error())
			messages = append(messages, msg)
	
			result.Code = 1
		} else {
			msg := util.SetMessage("INFO","Configuration [" + conf + "] create completed successfully")
			messages = append(messages, msg)
	
			result.Code = 0
		}
		result.Messages = messages
	
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

func AddProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var profileName string = params["profileName"]

	var result util.Result
	var messages []util.Message

	if profileName == "default" {
		msg := util.SetMessage("ERROR","Adding default profile not permitted!")
		messages = append(messages, msg)

		result.Code = 1	
		result.Messages = messages
		
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		dir := configDir + "/" + profileName
		err := util.CreateDir(dir,0755)
	
		if err != nil {
			msg := util.SetMessage("ERROR","Add profile [" + dir + "] failed!" + err.Error())
			messages = append(messages, msg)
	
			result.Code = 1
		} else {
			msg := util.SetMessage("INFO","Add profile [" + dir + "] completed successfully")
			messages = append(messages, msg)
	
			result.Code = 0
		}
		result.Messages = messages
	
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}	
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var profileName string = params["profileName"]

	var result util.Result
	var messages []util.Message

	dir := configDir + "/" + profileName

	if profileName == "default" {
		msg := util.SetMessage("ERROR","Deleting default profile not permitted!")
		messages = append(messages, msg)

		result.Code = 1	
		result.Messages = messages
		
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		err := util.RecursiveDirDelete(dir)

		if err != nil {
			msg := util.SetMessage("ERROR","Delete profile [" + dir + "] failed!" + err.Error())
			messages = append(messages, msg)
	
			result.Code = 1
		} else {
			msg := util.SetMessage("INFO","Delete profile [" + dir + "] completed successfully")
			messages = append(messages, msg)
	
			result.Code = 0
		}
		result.Messages = messages
	
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

func DeleteConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	var profileName string = params["profileName"]
	var configName string = params["configName"]

	var result util.Result
	var messages []util.Message

	if profileName == "default" {
		msg := util.SetMessage("ERROR","Deleting default configurations not permitted!")
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages	
		
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		configPath := configDir + "/" + profileName + "/" + configName + ".conf"
		err := os.Remove(configPath)
	
		if err != nil {
			msg := util.SetMessage("ERROR","Delete config [" + configPath + "] failed!" + err.Error())
			messages = append(messages, msg)
	
			result.Code = 1
		} else {
			msg := util.SetMessage("INFO","Delete config [" + configPath + "] completed successfully")
			messages = append(messages, msg)
	
			result.Code = 0
		}
		result.Messages = messages
	
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}	
}