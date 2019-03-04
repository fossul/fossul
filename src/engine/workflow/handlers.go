package main

import (
	"encoding/json"
	//    "io/ioutil"
	"engine/util"
	"net/http"
	//	"fmt"
	"log"
)

func GetStatusEndpoint(w http.ResponseWriter, r *http.Request) {
	var status = util.Status{Msg: "OK"}
	json.NewEncoder(w).Encode(status)
}

func StartBackupWorkflowEndpoint(w http.ResponseWriter, r *http.Request) {

	var config util.Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		log.Println(err)
	}
	defer r.Body.Close()
 
	res,err := json.Marshal(&config)
	if err != nil {
        log.Println(err)
    }
	
	
	log.Println("test", string(res), config.BackupRetentions)

	var quiesceResult util.Result
	quiesceResult = quiesce()

	var createBackupResult util.Result
	createBackupResult = createBackup()

	var unquiesceResult util.Result
	unquiesceResult = unquiesce()

	var deleteBackupResult util.Result
	deleteBackupResult = deleteBackup()

	var results []util.Result

	results = append(results, quiesceResult)
	results = append(results, createBackupResult)
	results = append(results, unquiesceResult)
	results = append(results, deleteBackupResult)

	_ = json.NewDecoder(r.Body).Decode(&results)
	json.NewEncoder(w).Encode(results)
}
