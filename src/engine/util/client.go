package util

import (
	"encoding/json"
	"log"
	"net/http"
	"bytes"
)

func GetWorkflowServiceStatus() Status {

	req, err := http.NewRequest("GET", "http://fossil-workflow:8000/status", nil)
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	var status Status

	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		log.Println(err)
	}

	return status

}

func GetAppServiceStatus() Status {

	req, err := http.NewRequest("GET", "http://fossil-app:8001/status", nil)
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	var status Status

	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		log.Println(err)
	}

	return status

}

func GetStorageServiceStatus() Status {

	req, err := http.NewRequest("GET", "http://fossil-storage:8002/status", nil)
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	var status Status

	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		log.Println(err)
	}

	return status

}

func StartBackupWorkflow(config Config) (result []Result) {


	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-workflow:8000/startBackupWorkflow", b)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	return result

}