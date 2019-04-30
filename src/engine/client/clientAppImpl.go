package client

import (
	"encoding/json"
	"fossil/src/engine/util"
	"log"
	"net/http"
	"bytes"
//	"strings"
)

func AppPluginList(pluginType string,config util.Config) []string {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-app:8001/pluginList/" + pluginType, b)
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

	var plugins []string

	if err := json.NewDecoder(resp.Body).Decode(&plugins); err != nil {
		log.Println(err)
	}

	return plugins

}

func AppPluginInfo(config util.Config, pluginName,pluginType string) (util.PluginInfoResult) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-app:8001/pluginInfo/" + pluginName + "/" + pluginType, b)
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

	var pluginInfoResult util.PluginInfoResult
	if err := json.NewDecoder(resp.Body).Decode(&pluginInfoResult); err != nil {
		log.Println(err)
	}

	return pluginInfoResult
}

func Discover(config util.Config) util.DiscoverResult {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-app:8001/discover", b)
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

	var discoverResult util.DiscoverResult

	if err := json.NewDecoder(resp.Body).Decode(&discoverResult); err != nil {
		log.Println(err)
	}

	return discoverResult
}

func Quiesce(config util.Config) util.Result {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-app:8001/quiesce", b)
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

	var result util.Result

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	return result
}

func Unquiesce(config util.Config) util.Result {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-app:8001/unquiesce", b)
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

	var result util.Result

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	return result
}