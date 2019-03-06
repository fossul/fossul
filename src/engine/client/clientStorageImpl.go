package client

import (
	"encoding/json"
	"engine/util"
	"log"
	"net/http"
)

func CreateBackup() util.Result {

	req, err := http.NewRequest("GET", "http://fossil-storage:8002/createBackup", nil)
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	var result util.Result

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	return result

}

func DeleteBackup() util.Result {

	req, err := http.NewRequest("GET", "http://fossil-storage:8002/deleteBackup", nil)
	if err != nil {
		log.Println("NewRequest: ", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
	}

	defer resp.Body.Close()

	var result util.Result

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	return result

}