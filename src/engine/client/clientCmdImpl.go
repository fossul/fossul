package client

import (
	"encoding/json"
	"engine/util"
	"log"
	"net/http"
)

func PreQuiesceCmd() util.Result {

	req, err := http.NewRequest("GET", "http://fossil-app:8001/preQuiesceCmd", nil)
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

func QuiesceCmd() util.Result {

	req, err := http.NewRequest("GET", "http://fossil-app:8001/quiesceCmd", nil)
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

func PostQuiesceCmd() util.Result {

	req, err := http.NewRequest("GET", "http://fossil-app:8001/postQuiesceCmd", nil)
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

func PreUnquiesceCmd() util.Result {

	req, err := http.NewRequest("GET", "http://fossil-app:8001/preUnquiesceCmd", nil)
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

func UnquiesceCmd() util.Result {

	req, err := http.NewRequest("GET", "http://fossil-app:8001/unquiesceCmd", nil)
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

func PostUnquiesceCmd() util.Result {

	req, err := http.NewRequest("GET", "http://fossil-app:8001/postUnquiesceCmd", nil)
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

func SendTrapSuccessCmd() util.Result {

	req, err := http.NewRequest("GET", "http://fossil-workflow:8000/sendTrapSuccessCmd", nil)
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

func SendTrapErrorCmd() util.Result {

	req, err := http.NewRequest("GET", "http://fossil-workflow:8000/sendTrapErrorCmd", nil)
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