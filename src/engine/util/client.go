package util

import (
	"encoding/json"
	"log"
	"net/http"
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