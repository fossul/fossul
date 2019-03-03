package main

import (
	"encoding/json"
	"engine/util"
	"net/http"
	"log"
//	"github.com/gorilla/mux"
)

func GetStatusEndpoint(w http.ResponseWriter, r *http.Request) {
	var status = util.Status{Msg: "OK"}
	json.NewEncoder(w).Encode(status)
}

func CreateQuiesceEndpoint(w http.ResponseWriter, r *http.Request) {
	//var result = Result{Code: 0, Stdout: "quiesce completed successfully", Stderr: "executed command xyz successfully"}
	//var config util.Config
	//_ = json.NewDecoder(r.Body).Decode(&config)
	//defer r.Body.Close()
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	var config util.Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		log.Println(err)
	}
 
	res,err := json.Marshal(&config)
	if err != nil {
        log.Println(err)
    }
	log.Println("test", string(res), config.BackupRententionDays)

//w.Write(res)
//	params := mux.Vars(r)
//	log.Println("test ", params["profile"])
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
