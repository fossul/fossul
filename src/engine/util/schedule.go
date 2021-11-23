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
package util

import (
	"encoding/json"
	"net/http"

	"github.com/robfig/cron/v3"
)

type JobScheduleResult struct {
	JobSchedules []JobSchedule `json:"jobSchedules,omitempty"`
	Result       Result        `json:"result,omitempty"`
}

type JobSchedule struct {
	CronId       cron.EntryID `json:"cronId"`
	CronSchedule string       `json:"cronSchedule"`
	ProfileName  string       `json:"profileName"`
	ConfigName   string       `json:"configName"`
	BackupPolicy string       `json:"backupPolicy"`
}

type CronSchedule struct {
	Value string `json:"value,omitempty"`
}

func GetCronSchedule(w http.ResponseWriter, r *http.Request) (CronSchedule, error) {

	var cronSchedule CronSchedule
	if err := json.NewDecoder(r.Body).Decode(&cronSchedule); err != nil {
		return cronSchedule, err
	}
	defer r.Body.Close()

	_, err := json.Marshal(&cronSchedule)
	if err != nil {
		return cronSchedule, err
	}

	return cronSchedule, nil
}

func GetCronScheduleId(w http.ResponseWriter, r *http.Request) (cron.EntryID, error) {

	var cronScheduleId cron.EntryID
	if err := json.NewDecoder(r.Body).Decode(&cronScheduleId); err != nil {
		return cronScheduleId, err
	}
	defer r.Body.Close()

	_, err := json.Marshal(&cronScheduleId)
	if err != nil {
		return cronScheduleId, err
	}

	return cronScheduleId, nil
}
