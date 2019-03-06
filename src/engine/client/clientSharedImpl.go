package client

import (
	"encoding/json"
	"engine/util"
	"log"
	"net/http"
	"bytes"
)

func PluginList(config util.Config) []string {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-app:8001/pluginList", b)
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

func PluginInfo(config util.Config, pluginName string) (util.Result, util.Plugin) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://fossil-app:8001/pluginInfo/" + pluginName, b)
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

	//unmarshall json response to plugin struct
	var plugin util.Plugin
	pluginBytes := []byte(result.Stdout)
    json.Unmarshal(pluginBytes, &plugin)

	return result, plugin
}

