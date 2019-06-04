package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fossul/src/engine/util"
	"net/http"
)

func AddSchedule(auth Auth, profileName, configName, policy, cronScheduleInput string) (util.Result, error) {
	var cronSchedule util.CronSchedule
	cronSchedule.Value = cronScheduleInput

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(cronSchedule)

	var result util.Result

	req, err := http.NewRequest("POST", "http://"+auth.ServerHostname+":"+auth.ServerPort+"/addSchedule/"+profileName+"/"+configName+"/"+policy, b)
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
