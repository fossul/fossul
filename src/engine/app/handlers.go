package main

import (
	"encoding/json"
	"engine/util"
	"net/http"
)

func GetStatus(w http.ResponseWriter, r *http.Request) {
	var status = util.Status{Msg: "OK"}
	json.NewEncoder(w).Encode(status)
}


func PreQuiesceCmd(w http.ResponseWriter, r *http.Request) {

	var result = util.SetResult(0, "pre quiesce cmd completed successfully", "executed command xyz successfully")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func QuiesceCmd(w http.ResponseWriter, r *http.Request) {

	var result = util.SetResult(0, "quiesce cmd completed successfully", "executed command xyz successfully")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func Quiesce(w http.ResponseWriter, r *http.Request) {

	var result util.Result
	result = util.ExecuteCommand("echo", "hello", "world")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func PostQuiesceCmd(w http.ResponseWriter, r *http.Request) {

	var result = util.SetResult(0, "post quiesce cmd completed successfully", "executed command xyz successfully")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func UnquiesceCmd(w http.ResponseWriter, r *http.Request) {

	var result = util.SetResult(0, "unquiesce cmd completed successfully", "executed command xyz successfully")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func PreUnquiesceCmd(w http.ResponseWriter, r *http.Request) {

	var result = util.SetResult(0, "pre unquiesce cmd completed successfully", "executed command xyz successfully")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func Unquiesce(w http.ResponseWriter, r *http.Request) {

	var result = util.SetResult(0, "unquiesce completed successfully", "executed command xyz successfully")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func PostUnquiesceCmd(w http.ResponseWriter, r *http.Request) {

	var result = util.SetResult(0, "post unquiesce cmd completed successfully", "executed command xyz successfully")
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}
