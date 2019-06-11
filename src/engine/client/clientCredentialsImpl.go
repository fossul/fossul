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

func AddAwsCredentials(auth Auth, awsCredentials util.AwsCredentials) (util.Result, error) {
	var result util.Result

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(awsCredentials)

	req, err := http.NewRequest("POST", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/addAwsCredentials", b)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return result, err
	}

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

func DeleteAwsCredentials(auth Auth) (util.Result, error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://"+auth.StorageHostname+":"+auth.StoragePort+"/deleteAwsCredentials", nil)
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		return result, err
	}

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
