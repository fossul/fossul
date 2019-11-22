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
	"github.com/fossul/fossul/src/engine/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

// GetConfig godoc
// @Description Get Configuration
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.ConfigResult
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /getConfig/{profileName}/{configName} [get]
func GetConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]

	var result util.Result
	var messages []util.Message
	var configResult util.ConfigResult

	conf := configDir + "/" + profileName + "/" + configName + "/" + configName + ".conf"

	if debug == "true" {
		log.Println("[DEBUG]", "Config path is "+conf)
	}
	config, err := util.ReadConfig(conf)
	printConfigDebug(config)

	if err != nil {
		msg := util.SetMessage("ERROR", "Couldn't get config! "+err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		configResult.Result = result

		_ = json.NewDecoder(r.Body).Decode(&configResult)
		json.NewEncoder(w).Encode(configResult)
	} else {
		result.Code = 0
		configResult.Result = result
		configResult.Config = config

		_ = json.NewDecoder(r.Body).Decode(&configResult)
		json.NewEncoder(w).Encode(configResult)
	}
}

// GetPluginConfig godoc
// @Description Get Plugin Configuration
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param pluginName path string true "name of plugin"
// @Accept  json
// @Produce  json
// @Success 200 {map} string
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /getPluginConfig/{profileName}/{configName}/{pluginName} [get]
func GetPluginConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var pluginName string = params["pluginName"]

	var result util.Result
	var messages []util.Message
	var configMapResult util.ConfigMapResult

	conf := configDir + "/" + profileName + "/" + configName + "/" + pluginName + ".conf"

	if debug == "true" {
		log.Println("[DEBUG]", "Plugin config path is "+conf)
	}

	configMap, err := util.ReadConfigToMap(conf)
	printConfigMapDebug(configMap)

	if err != nil {
		msg := util.SetMessage("ERROR", "Couldn't get config! "+err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		configMapResult.Result = result

		_ = json.NewDecoder(r.Body).Decode(&configMapResult)
		json.NewEncoder(w).Encode(configMapResult)
	} else {
		result.Code = 0
		configMapResult.Result = result
		configMapResult.ConfigMap = configMap

		_ = json.NewDecoder(r.Body).Decode(&configMapResult)
		json.NewEncoder(w).Encode(configMapResult)
	}
}

// GetDefaultConfig godoc
// @Description Get Default Configuration
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Config
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /getDefaultConfig [get]
func GetDefaultConfig(w http.ResponseWriter, r *http.Request) {

	var result util.Result
	var messages []util.Message
	var configResult util.ConfigResult

	conf := configDir + "/" + "default" + "/" + "default" + "/" + "default.conf"

	if debug == "true" {
		log.Println("[DEBUG]", "Default config path is "+conf)
	}

	config, err := util.ReadConfig(conf)
	printConfigDebug(config)

	if err != nil {
		msg := util.SetMessage("ERROR", "Couldn't get config! "+err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages
		configResult.Result = result

		_ = json.NewDecoder(r.Body).Decode(&configResult)
		json.NewEncoder(w).Encode(configResult)
	} else {
		result.Code = 0
		configResult.Result = result

		configResult.Config = config

		_ = json.NewDecoder(r.Body).Decode(&configResult)
		json.NewEncoder(w).Encode(configResult)
	}
}

// GetDefaultPluginConfig godoc
// @Description Get Default Plugin Configuration
// @Param pluginName path string true "name of plugin"
// @Accept  json
// @Produce  json
// @Success 200 {map} string
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /getDefaultPluginConfig/{pluginName} [get]
func GetDefaultPluginConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var pluginName string = params["pluginName"]

	conf := configDir + "/" + "default" + "/" + "default" + "/" + pluginName + ".conf"

	var result util.Result
	var messages []util.Message
	var configMapResult util.ConfigMapResult

	if debug == "true" {
		log.Println("[DEBUG]", "Config path is "+conf)
	}

	configMap, err := util.ReadConfigToMap(conf)
	printConfigMapDebug(configMap)

	if err != nil {
		msg := util.SetMessage("ERROR", "Couldn't get config! "+err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		configMapResult.Result = result

		_ = json.NewDecoder(r.Body).Decode(&configMapResult)
		json.NewEncoder(w).Encode(configMapResult)
	} else {
		result.Code = 0
		configMapResult.Result = result
		configMapResult.ConfigMap = configMap

		_ = json.NewDecoder(r.Body).Decode(&configMapResult)
		json.NewEncoder(w).Encode(configMapResult)
	}
}

// AddConfig godoc
// @Description Add Configuration
// @Param config body util.Config true "config struct"
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /addConfig/{profileName}/{configName} [post]
func AddConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]

	var result util.Result
	var messages []util.Message

	config, err := util.GetConfig(w, r)
	printConfigDebug(config)

	if err != nil {
		msg := util.SetMessage("ERROR", "Couldn't get configuration. "+err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	dir := configDir + "/" + profileName + "/" + configName
	err = util.CreateDir(dir, 0755)
	if err != nil {
		msg := util.SetMessage("ERROR", "Couldn't create config directory! "+err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}
	conf := dir + "/" + configName + ".conf"

	if util.ExistsPath(dir) != true {
		msg := util.SetMessage("ERROR", "Add configuration ["+conf+"] failed! Profile ["+dir+"] does not exist.")
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		//err := util.WriteGob(conf,config)
		err := util.WriteConfig(conf, config)

		if err != nil {
			msg := util.SetMessage("ERROR", "Add configuration ["+conf+"] failed!"+err.Error())
			messages = append(messages, msg)

			result.Code = 1
		} else {
			msg := util.SetMessage("INFO", "Configuration ["+conf+"] create completed successfully")
			messages = append(messages, msg)

			result.Code = 0
		}
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// AddPluginConfig godoc
// @Description Add Plugin Configuration
// @Param config body util.PluginConfigMap true "config map"
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param pluginName path string true "name of plugin"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /addPluginConfig/{profileName}/{configName}/{pluginName} [post]
func AddPluginConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var pluginName string = params["pluginName"]

	var result util.Result
	var messages []util.Message

	configMap, err := util.GetPluginConfig(w, r)
	printConfigMapDebug(configMap)

	if err != nil {
		msg := util.SetMessage("ERROR", "Couldn't get configuration. "+err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	dir := configDir + "/" + profileName + "/" + configName
	err = util.CreateDir(dir, 0755)
	if err != nil {
		msg := util.SetMessage("ERROR", "Couldn't create config directory! "+err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	conf := dir + "/" + pluginName + ".conf"

	if util.ExistsPath(dir) != true {
		msg := util.SetMessage("ERROR", "Add configuration ["+conf+"] failed! Profile ["+dir+"] does not exist.")
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		//err := util.WriteGob(conf,config)
		err := util.WritePluginConfig(conf, configMap)

		if err != nil {
			msg := util.SetMessage("ERROR", "Add configuration ["+conf+"] failed!"+err.Error())
			messages = append(messages, msg)

			result.Code = 1
		} else {
			msg := util.SetMessage("INFO", "Configuration ["+conf+"] create completed successfully")
			messages = append(messages, msg)

			result.Code = 0
		}
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// DeletePluginConfig godoc
// @Description Add Plugin Configuration
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param pluginName path string true "name of plugin"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /deletePluginConfig/{profileName}/{configName}/{pluginName} [get]
func DeletePluginConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var pluginName string = params["pluginName"]

	var result util.Result
	var messages []util.Message

	if profileName == "default" {
		msg := util.SetMessage("ERROR", "Deleting default config not permitted!")
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		configPath := configDir + "/" + profileName + "/" + configName + "/" + pluginName + ".conf"
		err := os.Remove(configPath)

		if err != nil {
			msg := util.SetMessage("ERROR", "Delete plugin config ["+configPath+"] failed!"+err.Error())
			messages = append(messages, msg)

			result.Code = 1
		} else {
			msg := util.SetMessage("INFO", "Delete plugin config ["+configPath+"] completed successfully")
			messages = append(messages, msg)

			result.Code = 0
		}
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// AddProfile godoc
// @Description Add Profile
// @Param profileName path string true "name of profile"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /addProfile/{profileName} [get]
func AddProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]

	var result util.Result
	var messages []util.Message

	if profileName == "default" {
		msg := util.SetMessage("ERROR", "Adding default profile not permitted!")
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		dir := configDir + "/" + profileName
		err := util.CreateDir(dir, 0755)

		if err != nil {
			msg := util.SetMessage("ERROR", "Add profile ["+dir+"] failed!"+err.Error())
			messages = append(messages, msg)

			result.Code = 1
		} else {
			msg := util.SetMessage("INFO", "Add profile ["+dir+"] completed successfully")
			messages = append(messages, msg)

			result.Code = 0
		}
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// DeleteProfile godoc
// @Description Delete Profile Including Configurations (destructive)
// @Param profileName path string true "name of profile"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /deleteProfile/{profileName} [get]
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]

	var result util.Result
	var messages []util.Message

	dir := configDir + "/" + profileName

	if profileName == "default" {
		msg := util.SetMessage("ERROR", "Deleting default profile not permitted!")
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		err := util.RecursiveDirDelete(dir)

		if err != nil {
			msg := util.SetMessage("ERROR", "Delete profile ["+dir+"] failed!"+err.Error())
			messages = append(messages, msg)

			result.Code = 1
		} else {
			msg := util.SetMessage("INFO", "Delete profile ["+dir+"] completed successfully")
			messages = append(messages, msg)

			result.Code = 0
		}
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// DeleteConfigDir godoc
// @Description Delete Entire Configuration (destructive)
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /deleteConfigDir/{profileName}/{configName} [get]
func DeleteConfigDir(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]

	var result util.Result
	var messages []util.Message

	dir := configDir + "/" + profileName + "/" + configName

	if profileName == "default" {
		msg := util.SetMessage("ERROR", "Deleting default config not permitted!")
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		err := util.RecursiveDirDelete(dir)

		if err != nil {
			msg := util.SetMessage("ERROR", "Delete config dir ["+dir+"] failed!"+err.Error())
			messages = append(messages, msg)

			result.Code = 1
		} else {
			msg := util.SetMessage("INFO", "Delete config dir ["+dir+"] completed successfully")
			messages = append(messages, msg)

			result.Code = 0
		}
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// DeleteConfig godoc
// @Description Delete Configuration
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /deleteConfig/{profileName}/{configName} [get]
func DeleteConfig(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]

	var result util.Result
	var messages []util.Message

	if profileName == "default" {
		msg := util.SetMessage("ERROR", "Deleting default configurations not permitted!")
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		configPath := configDir + "/" + profileName + "/" + configName + "/" + configName + ".conf"
		err := os.Remove(configPath)

		if err != nil {
			msg := util.SetMessage("ERROR", "Delete config ["+configPath+"] failed!"+err.Error())
			messages = append(messages, msg)

			result.Code = 1
		} else {
			msg := util.SetMessage("INFO", "Delete config ["+configPath+"] completed successfully")
			messages = append(messages, msg)

			result.Code = 0
		}
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// ListProfiles godoc
// @Description List Profiles
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /listProfiles [get]
func ListProfiles(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var messages []util.Message

	profiles, err := util.DirectoryList(configDir)
	if err != nil {
		msg := util.SetMessage("ERROR", "Profile list failed! "+err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		result.Code = 0
		result.Data = profiles

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}

// ListConfigs godoc
// @Description List Configurations
// @Param profileName path string true "name of profile"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /listConfigs/{profileName} [get]
func ListConfigs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]

	var result util.Result
	var messages []util.Message
	var configs []string

	configs, err := util.DirectoryList(configDir + "/" + profileName)
	if err != nil {
		msg := util.SetMessage("ERROR", "Config list failed! "+err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
		return
	}

	result.Code = 0
	result.Data = configs

	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

// ListPluginConfigs godoc
// @Description List Plugin Configuration
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /listPluginConfigs/{profileName}/{configName} [get]
func ListPluginConfigs(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]

	var result util.Result
	var messages []util.Message

	plugins, err := util.PluginList(configDir+"/"+profileName+"/"+configName, configName)
	if err != nil {
		msg := util.SetMessage("ERROR", "Plugin list failed! "+err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	} else {
		result.Code = 0
		result.Data = plugins

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	}
}
