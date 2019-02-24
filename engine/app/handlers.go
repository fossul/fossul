package main

import (
	"encoding/json"
	"net/http"
)
func GetStatusEndpoint(w http.ResponseWriter, r *http.Request) {
	var status = Status{Msg: "OK"}
	json.NewEncoder(w).Encode(status)
}

func CreateQuiesceEndpoint(w http.ResponseWriter, r *http.Request) {
	var result = Result{Code: 0, Stdout: "quiesce completed successfully", Stderr: "executed command xyz successfully"}
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func CreateUnquiesceEndpoint(w http.ResponseWriter, r *http.Request) {
	var result = Result{Code: 0, Stdout: "unquiesce completed successfully", Stderr: "executed command xyz successfully"}
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}