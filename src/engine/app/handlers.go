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

func CreateQuiesceEndpoint(w http.ResponseWriter, r *http.Request) {

	var result util.Result
	result = util.ExecuteCommand("echo", "hello", "world")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func CreateUnquiesceEndpoint(w http.ResponseWriter, r *http.Request) {
	var result = util.Result{Code: 0, Stdout: "unquiesce completed successfully", Stderr: "executed command xyz successfully"}
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}
