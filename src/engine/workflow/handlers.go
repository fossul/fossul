package main

import (
	"encoding/json"
	//    "io/ioutil"
	"engine/util"
	"net/http"
	//	"fmt"
)

func GetStatusEndpoint(w http.ResponseWriter, r *http.Request) {
	var status = util.Status{Msg: "OK"}
	json.NewEncoder(w).Encode(status)
}

func CreateBackupWorkflowEndpoint(w http.ResponseWriter, r *http.Request) {

	/*  OLD CODE
	    req, err := http.NewRequest("GET", "http://localhost:8001/quiesce", nil)
		if err != nil {
			log.Fatal("NewRequest: ", err)
			return
		}

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Do: ", err)
			return
		}

		defer resp.Body.Close()

		var quiesce Result

		if err := json.NewDecoder(resp.Body).Decode(&quiesce); err != nil {
			log.Println(err)
		}

	//	response, err := http.Get("http://localhost:8001/quiesce")
	//	data, _ := ioutil.ReadAll(response.Body)
	*/
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
