package main

import (
    "encoding/json"
//    "io/ioutil"
	"net/http"
	"log"
//	"fmt"
)

func GetStatusEndpoint(w http.ResponseWriter, r *http.Request) {
	var status = Status{Msg: "OK"}
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
	var quiesceResult Result
	quiesceResult = quiesce()

	var createBackupResult Result
	createBackupResult = createBackup()

	var unquiesceResult Result
	unquiesceResult = unquiesce()

	var deleteBackupResult Result
	deleteBackupResult = deleteBackup()

	var results []Result

	results = append(results, quiesceResult)
	results = append(results, createBackupResult)
	results = append(results, unquiesceResult)
	results = append(results, deleteBackupResult)

	_ = json.NewDecoder(r.Body).Decode(&results)
	json.NewEncoder(w).Encode(results)
}

func quiesce () Result {

	req, err := http.NewRequest("GET", "http://fossil-app:8001/quiesce", nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}

	defer resp.Body.Close()

	var result Result

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	return result

}

func createBackup () Result {

	req, err := http.NewRequest("GET", "http://fossil-storage:8002/createBackup", nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}

	defer resp.Body.Close()

	var result Result

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	return result

}

func deleteBackup () Result {

	req, err := http.NewRequest("GET", "http://fossil-storage:8002/deleteBackup", nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}

	defer resp.Body.Close()

	var result Result

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	return result

}

func unquiesce () Result {

	req, err := http.NewRequest("GET", "http://fossil-app:8001/unquiesce", nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}

	defer resp.Body.Close()

	var result Result

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	return result

}

