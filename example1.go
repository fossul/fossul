package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Result struct {
	Code    int      `json:"code"`
	Stdout  string   `json:"stdout"`
	Stderr  string   `json:"stderr"`
}

type Status struct {
	Msg string `json:"msg"`
}

func GetStatusEndpoint(w http.ResponseWriter, r *http.Request) {
	var status = Status{Msg: "OK"}
	json.NewEncoder(w).Encode(status)
}

func CreateQuiesceEndpoint(w http.ResponseWriter, r *http.Request) {
	//var result []Result
	var result = Result{Code: 0, Stdout: "quiesce completed successfully", Stderr: "executed command xyz successfully"}
	//params := mux.Vars(r)
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func CreateUnquiesceEndpoint(w http.ResponseWriter, r *http.Request) {
	var result = Result{Code: 0, Stdout: "quiesce completed successfully", Stderr: "executed command xyz successfully"}
	//result = Result{Code: 0, Stdout: "unquiesce completed successfully", Stderr: "executed command xyz successfully"}
	//params := mux.Vars(r)
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}


func main() {
	router := mux.NewRouter()
	router.HandleFunc("/status", GetStatusEndpoint).Methods("GET")
	router.HandleFunc("/quiesce", CreateQuiesceEndpoint).Methods("POST")
	router.HandleFunc("/unquiesce", CreateUnquiesceEndpoint).Methods("POST")
	log.Fatal(http.ListenAndServe(":8001", router))
}