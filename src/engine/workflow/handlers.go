package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"engine/util"
	"engine/client"
	"net/http"
	"strings"
	"log"
)

func GetStatus(w http.ResponseWriter, r *http.Request) {
	var status = util.Status{Msg: "OK"}
	json.NewEncoder(w).Encode(status)
}

func StartBackupWorkflow(w http.ResponseWriter, r *http.Request) {

	var config util.Config = util.GetConfig(w,r)

	var sendTrapErrorCmdResult util.Result
	var results []util.Result

	commentResult := util.SetResultMessage(0,"COMMENT","Welcome to Fossil Backup Framework, Performing Backup Workflow")
	results = append(results, commentResult)


	commentResult = util.SetResultMessage(0,"COMMENT","Performing Application Quiesce")
	results = append(results, commentResult)

	if config.PreAppQuiesceCmd != "" {
		preQuiesceCmdResult := client.PreQuiesceCmd(config)
		results = append(results, preQuiesceCmdResult)

		if preQuiesceCmdResult.Code != 0 {
			sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
			sendError(w,r,results)
		}
	}

	if config.AppQuiesceCmd != "" {
		quiesceCmdResult := client.QuiesceCmd(config)	
		results = append(results, quiesceCmdResult)

		if quiesceCmdResult.Code != 0 {
			sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
			sendError(w,r,results)
		}
	}	
	
	if config.AppPlugin != "" {
		quiesceResult := client.Quiesce(config)
		results = append(results, quiesceResult)	

		if quiesceResult.Code != 0 {
			sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
			sendError(w,r,results)
		}
	}	

	if config.PostAppQuiesceCmd != "" {
		postQuiesceCmdResult := client.PostQuiesceCmd(config)
		results = append(results, postQuiesceCmdResult)

		if postQuiesceCmdResult.Code != 0 {
			sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
			sendError(w,r,results)
		}
	}	

	commentResult = util.SetResultMessage(0,"COMMENT","Performing Backup")
	results = append(results, commentResult)

	if config.BackupCreateCmd != "" {
		backupCreateCmdResult := client.BackupCreateCmd(config)	
		results = append(results, backupCreateCmdResult)

		if backupCreateCmdResult.Code != 0 {
			sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
			sendError(w,r,results)
		}
	}	

	if config.StoragePlugin != "" {	
		backupResult := client.Backup(config)
		results = append(results, backupResult)

		if backupResult.Code != 0 {
			sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
			sendError(w,r,results)
		}
	}	

	commentResult = util.SetResultMessage(0,"COMMENT","Performing Application Unquiesce")
	results = append(results, commentResult)

	if config.PreAppUnquiesceCmd != "" {
		preUnquiesceCmdResult := client.PreUnquiesceCmd(config)
		results = append(results, preUnquiesceCmdResult)

		if preUnquiesceCmdResult.Code != 0 {
			sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
			sendError(w,r,results)
		}
	}	

	if config.AppUnquiesceCmd != "" {
		unquiesceCmdResult := client.UnquiesceCmd(config)	
		results = append(results, unquiesceCmdResult)

		if unquiesceCmdResult.Code != 0 {
			sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
			sendError(w,r,results)
		}
	}	

	if config.AppPlugin != "" {
		unquiesceResult := client.Unquiesce(config)
		results = append(results, unquiesceResult)

		if unquiesceResult.Code != 0 {
			sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
			sendError(w,r,results)
		}
	}	

	if config.PostAppUnquiesceCmd != "" {
		postUnquiesceCmdResult := client.PostUnquiesceCmd(config)
		results = append(results, postUnquiesceCmdResult)

		if postUnquiesceCmdResult.Code != 0 {
			sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
			sendError(w,r,results)	
		}
	}	

	commentResult = util.SetResultMessage(0,"COMMENT","Performing Backup Retention")
	results = append(results, commentResult)

	if config.BackupDeleteCmd != "" {
		backupDeleteCmdResult := client.BackupDeleteCmd(config)	
		results = append(results, backupDeleteCmdResult)

		if backupDeleteCmdResult.Code != 0 {
			sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
			sendError(w,r,results)
		}
	}	


	if config.StoragePlugin != "" {	
		backupDeleteResult := client.BackupDelete(config)
		results = append(results, backupDeleteResult)

		if backupDeleteResult.Code != 0 {
			sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
			sendError(w,r,results)
		}
	}	

	if config.SendTrapSuccessCmd != "" {
		sendTrapSuccessCmdResult := client.SendTrapSuccessCmd(config)	
		results = append(results, sendTrapSuccessCmdResult)

		if sendTrapSuccessCmdResult.Code != 0 {
			sendTrapErrorCmdResult = client.SendTrapErrorCmd(config)
			sendError(w,r,results)
		}
	}	

	commentResult = util.SetResultMessage(0,"INFO","Backup Completed Successfully")
	results = append(results, commentResult)
	
	results = append(results, sendTrapErrorCmdResult)

	_ = json.NewDecoder(r.Body).Decode(&results)
	json.NewEncoder(w).Encode(results)
}

func sendError(w http.ResponseWriter, r *http.Request, results []util.Result) {
	_ = json.NewDecoder(r.Body).Decode(&results)
	json.NewEncoder(w).Encode(results)
}

func SendTrapSuccessCmd(w http.ResponseWriter, r *http.Request) {
	var result util.Result

	var config util.Config = util.GetConfig(w,r)

	if config.SendTrapSuccessCmd != "" {
		args := strings.Split(config.SendTrapSuccessCmd, ",")
		message := util.SetMessage("INFO", "Performing send trap success command")

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
		message := util.SetMessage("INFO", "Performing send trap error command")

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