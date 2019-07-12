/*
Copyright 2019 The Fossul Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fossul/src/engine/util"
	"net/http"
)

func PreQuiesceCmd(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.AppHostname+":"+auth.AppPort+"/preQuiesceCmd", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func QuiesceCmd(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.AppHostname+":"+auth.AppPort+"/quiesceCmd", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func PostQuiesceCmd(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.AppHostname+":"+auth.AppPort+"/postQuiesceCmd", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func PreUnquiesceCmd(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.AppHostname+":"+auth.AppPort+"/preUnquiesceCmd", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func UnquiesceCmd(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.AppHostname+":"+auth.AppPort+"/unquiesceCmd", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func BackupCreateCmd(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/backupCreateCmd", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func BackupDeleteCmd(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/backupDeleteCmd", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func PostUnquiesceCmd(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.AppHostname+":"+auth.AppPort+"/postUnquiesceCmd", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func SendTrapSuccessCmd(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/sendTrapSuccessCmd", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func SendTrapErrorCmd(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/sendTrapErrorCmd", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func ArchiveCreateCmd(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/archiveCreateCmd", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func ArchiveDeleteCmd(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/archiveDeleteCmd", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func PreAppRestoreCmd(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.AppHostname+":"+auth.AppPort+"/preAppRestoreCmd", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func RestoreCmd(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/restoreCmd", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}

func PostAppRestoreCmd(auth Auth, config util.Config) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(config)

	req, err := http.NewRequest("POST", "http://"+auth.AppHostname+":"+auth.AppPort+"/postAppRestoreCmd", b)
	if err != nil {
		return result, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return result, err
		}
	} else {
		return result, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return result, nil
}
