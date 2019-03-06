package main

import (
	"encoding/json"
	"engine/util"
	"engine/client"
	"net/http"
)

func GetStatus(w http.ResponseWriter, r *http.Request) {
	var status = util.Status{Msg: "OK"}
	json.NewEncoder(w).Encode(status)
}

func StartBackupWorkflow(w http.ResponseWriter, r *http.Request) {

	var config util.Config = util.GetConfig(w,r)

	var sendTrapErrorCmdResult util.Result

	var preQuiesceCmdResult util.Result
	preQuiesceCmdResult = client.PreQuiesceCmd()
	if preQuiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd()
	}

	var quiesceCmdResult util.Result
	quiesceCmdResult = client.QuiesceCmd()
	if quiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd()
	}
	
	var quiesceResult util.Result
	quiesceResult = client.Quiesce(config)
	if quiesceResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd()
	}

	var postQuiesceCmdResult util.Result
	postQuiesceCmdResult = client.PostQuiesceCmd()
	if postQuiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd()
	}

	var createBackupResult util.Result
	createBackupResult = client.CreateBackup()
	if createBackupResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd()
	}

	var preUnquiesceCmdResult util.Result
	preUnquiesceCmdResult = client.PreUnquiesceCmd()
	if preUnquiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd()
	}

	var unquiesceCmdResult util.Result
	unquiesceCmdResult = client.UnquiesceCmd()
	if unquiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd()
	}

	var unquiesceResult util.Result
	unquiesceResult = client.Unquiesce(config)
	if unquiesceResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd()
	}

	var postUnquiesceCmdResult util.Result
	postUnquiesceCmdResult = client.PostUnquiesceCmd()
	if postUnquiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd()
	}

	var deleteBackupResult util.Result
	deleteBackupResult = client.DeleteBackup()
	if deleteBackupResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd()
	}

	var sendTrapSuccessCmdResult util.Result
	sendTrapSuccessCmdResult = client.SendTrapSuccessCmd()
	if sendTrapSuccessCmdResult.Code != 0 {
		sendTrapErrorCmdResult = client.SendTrapErrorCmd()
	}

	var results []util.Result

	if (util.Result{}) != preQuiesceCmdResult {
		results = append(results, preQuiesceCmdResult)
	}
	if (util.Result{}) != quiesceCmdResult {	
		results = append(results, quiesceCmdResult)
	}
	if (util.Result{}) != quiesceResult {	
		results = append(results, quiesceResult)
	}	
	if (util.Result{}) != postQuiesceCmdResult {
		results = append(results, postQuiesceCmdResult)
	}
	if (util.Result{}) != createBackupResult {	
		results = append(results, createBackupResult)
	}
	if (util.Result{}) != preUnquiesceCmdResult {
		results = append(results, preUnquiesceCmdResult)
	}
	if (util.Result{}) != unquiesceCmdResult {	
		results = append(results, unquiesceCmdResult)
	}
	if (util.Result{}) != unquiesceResult {
		results = append(results, unquiesceResult)
	}
	if (util.Result{}) != postUnquiesceCmdResult {
		results = append(results, postUnquiesceCmdResult)
	}
	if (util.Result{}) != deleteBackupResult {	
		results = append(results, deleteBackupResult)
	}
	if (util.Result{}) != sendTrapSuccessCmdResult {	
		results = append(results, sendTrapSuccessCmdResult)
	}
	if (util.Result{}) != sendTrapErrorCmdResult {	
		results = append(results, sendTrapErrorCmdResult)
	}

	_ = json.NewDecoder(r.Body).Decode(&results)
	json.NewEncoder(w).Encode(results)
}

func SendTrapSuccessCmd(w http.ResponseWriter, r *http.Request) {

	var result = util.SetResult(0, "send trap success cmd completed successfully", "executed command xyz successfully")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func SendTrapErrorCmd(w http.ResponseWriter, r *http.Request) {

	var result = util.SetResult(0, "send trap error cmd completed successfully", "executed command xyz successfully")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}