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
	"net/http"

	"github.com/fossul/fossul/src/engine/util"
)

func AddSchedule(auth Auth, profileName, configName, policy, cronScheduleInput string) (util.Result, error) {
	var cronSchedule util.CronSchedule
	cronSchedule.Value = cronScheduleInput

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(cronSchedule)

	var result util.Result

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/addSchedule/"+profileName+"/"+configName+"/"+policy, b)
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

func DeleteSchedule(auth Auth, profileName, configName, policy string) (util.Result, error) {
	var result util.Result

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/deleteSchedule/"+profileName+"/"+configName+"/"+policy, nil)
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

func ListSchedules(auth Auth) (util.JobScheduleResult, error) {
	var jobScheduleResult util.JobScheduleResult

	req, err := http.NewRequest("GET", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/listSchedules", nil)
	if err != nil {
		return jobScheduleResult, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(auth.Username, auth.Password)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return jobScheduleResult, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		if err := json.NewDecoder(resp.Body).Decode(&jobScheduleResult); err != nil {
			return jobScheduleResult, err
		}
	} else {
		return jobScheduleResult, errors.New("Http Status Error [" + resp.Status + "]")
	}

	return jobScheduleResult, nil
}
