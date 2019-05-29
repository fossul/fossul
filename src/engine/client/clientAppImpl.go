package client

import (
	"encoding/json"
	"fossil/src/engine/util"
	"net/http"
	"bytes"
	"errors"
//	"strings"
)

func AppPluginList(auth Auth,pluginType string) ([]string,error) {
	var plugins []string

	req, err := http.NewRequest("GET", "http://" + auth.AppHostname + ":" + auth.AppPort + "/pluginList/" + pluginType, nil)
	if err != nil {
		return plugins,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return plugins,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&plugins); err != nil {
			return plugins,err
		}
	} else {
		return plugins,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return plugins,nil

}

func AppPluginInfo(auth Auth,config util.Config, pluginName,pluginType string) (util.PluginInfoResult,error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	var pluginInfoResult util.PluginInfoResult

	req, err := http.NewRequest("POST", "http://" + auth.AppHostname + ":" + auth.AppPort + "/pluginInfo/" + pluginName + "/" + pluginType, b)
	if err != nil {
		return pluginInfoResult,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return pluginInfoResult,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&pluginInfoResult); err != nil {
			return pluginInfoResult,err
		}
	} else {
		return pluginInfoResult,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return pluginInfoResult,nil
}

func Discover(auth Auth,config util.Config) (util.DiscoverResult,error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	var discoverResult util.DiscoverResult

	req, err := http.NewRequest("POST", "http://" + auth.AppHostname + ":" + auth.AppPort + "/discover", b)
	if err != nil {
		return discoverResult,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return discoverResult,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&discoverResult); err != nil {
			return discoverResult,err
		}
	} else {
		return discoverResult,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return discoverResult,nil
}

func Quiesce(auth Auth,config util.Config) (util.Result,error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	var result util.Result

	req, err := http.NewRequest("POST", "http://" + auth.AppHostname + ":" + auth.AppPort + "/quiesce", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func Unquiesce(auth Auth,config util.Config) (util.Result,error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	var result util.Result

	req, err := http.NewRequest("POST", "http://" + auth.AppHostname + ":" + auth.AppPort + "/unquiesce", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()


	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func PreRestore(auth Auth,config util.Config) (util.Result,error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	var result util.Result

	req, err := http.NewRequest("POST", "http://" + auth.AppHostname + ":" + auth.AppPort + "/preRestore", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}

func PostRestore(auth Auth,config util.Config) (util.Result,error) {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	var result util.Result

	req, err := http.NewRequest("POST", "http://" + auth.AppHostname + ":" + auth.AppPort + "/postRestore", b)
	if err != nil {
		return result,err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result,err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result,err
		}
	} else {
		return result,errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result,nil
}
