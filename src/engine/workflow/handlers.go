package main

import (
	"encoding/json"
	"engine/util"
	"net/http"
	"log"
)

func GetStatus(w http.ResponseWriter, r *http.Request) {
	var status = util.Status{Msg: "OK"}
	json.NewEncoder(w).Encode(status)
}

func StartBackupWorkflow(w http.ResponseWriter, r *http.Request) {

	var config util.Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		log.Println(err)
	}
	defer r.Body.Close()
 
	res,err := json.Marshal(&config)
	if err != nil {
        log.Println(err)
    }
	
	
	log.Println("DEBUG", string(res), config.BackupRetentions)

	var sendTrapErrorCmdResult util.Result

	var preQuiesceCmdResult util.Result
	preQuiesceCmdResult = preQuiesceCmd()
	if preQuiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = sendTrapErrorCmd()
	}

	var quiesceCmdResult util.Result
	quiesceCmdResult = quiesceCmd()
	if quiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = sendTrapErrorCmd()
	}
	
	var quiesceResult util.Result
	quiesceResult = quiesce()
	if quiesceResult.Code != 0 {
		sendTrapErrorCmdResult = sendTrapErrorCmd()
	}

	var postQuiesceCmdResult util.Result
	postQuiesceCmdResult = postQuiesceCmd()
	if postQuiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = sendTrapErrorCmd()
	}

	var createBackupResult util.Result
	createBackupResult = createBackup()
	if createBackupResult.Code != 0 {
		sendTrapErrorCmdResult = sendTrapErrorCmd()
	}

	var preUnquiesceCmdResult util.Result
	preUnquiesceCmdResult = preUnquiesceCmd()
	if preUnquiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = sendTrapErrorCmd()
	}

	var unquiesceCmdResult util.Result
	unquiesceCmdResult = unquiesceCmd()
	if unquiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = sendTrapErrorCmd()
	}

	var unquiesceResult util.Result
	unquiesceResult = unquiesce()
	if unquiesceResult.Code != 0 {
		sendTrapErrorCmdResult = sendTrapErrorCmd()
	}

	var postUnquiesceCmdResult util.Result
	postUnquiesceCmdResult = postUnquiesceCmd()
	if postUnquiesceCmdResult.Code != 0 {
		sendTrapErrorCmdResult = sendTrapErrorCmd()
	}

	var deleteBackupResult util.Result
	deleteBackupResult = deleteBackup()
	if deleteBackupResult.Code != 0 {
		sendTrapErrorCmdResult = sendTrapErrorCmd()
	}

	var sendTrapSuccessCmdResult util.Result
	sendTrapSuccessCmdResult = sendTrapSuccessCmd()
	if sendTrapSuccessCmdResult.Code != 0 {
		sendTrapErrorCmdResult = sendTrapErrorCmd()
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