package main

import (
	"encoding/json"
	"engine/util"
	"net/http"
)

func GetStatusEndpoint(w http.ResponseWriter, r *http.Request) {
	var status = util.Status{Msg: "OK"}
	json.NewEncoder(w).Encode(status)
}

func CreateBackupEndpoint(w http.ResponseWriter, r *http.Request) {
	var result = util.Result{Code: 0, Stdout: "backup create completed successfully", Stderr: "executed command xyz successfully"}
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func DeleteBackupEndpoint(w http.ResponseWriter, r *http.Request) {
	var result = util.Result{Code: 0, Stdout: "backup delete completed successfully", Stderr: "executed command xyz successfully"}
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}
