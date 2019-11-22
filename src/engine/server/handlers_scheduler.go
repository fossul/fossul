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
package main

import (
	"encoding/json"
	"github.com/fossul/fossul/src/engine/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// AddSchedule godoc
// @Description Add job schedule
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param policy path string true "policy name"
// @Param cronSchedule body util.CronSchedule true "value: min,hour,dayOfMonth,month,dayOfWeek"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /addSchedule/{profileName}/{configName}/{policy} [post]
func AddSchedule(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var policy string = params["policy"]

	var result util.Result
	var messages []util.Message

	cronSchedule, err := util.GetCronSchedule(w, r)
	if err != nil {
		msg := util.SetMessage("ERROR", "Couldn't get cron schedule! "+err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	id, err := AddCronSchedule(profileName, configName, policy, cronSchedule.Value)
	if err != nil {
		msg := util.SetMessage("ERROR", "Add schedule failed! "+err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	if debug == "true" {
		log.Println("[DEBUG] Added schedule with id", id)
	}

	msg := util.SetMessage("INFO", "Add schedule completed successfully")
	messages = append(messages, msg)

	result.Code = 0
	result.Messages = messages

	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

// DeleteSchedule godoc
// @Description Delete job schedule
// @Param profileName path string true "name of profile"
// @Param configName path string true "name of config"
// @Param policy path string true "policy name"
// @Accept  json
// @Produce  json
// @Success 200 {object} util.Result
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /deleteSchedule/{profileName}/{configName}/{policy} [get]
func DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var policy string = params["policy"]

	var result util.Result
	var messages []util.Message

	err := DeleteCronSchedule(profileName, configName, policy)
	if err != nil {
		msg := util.SetMessage("ERROR", "Delete schedule failed! "+err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)

		return
	}

	msg := util.SetMessage("INFO", "Delete schedule completed successfully")
	messages = append(messages, msg)

	result.Code = 0
	result.Messages = messages

	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

// ListSchedules godoc
// @Description List job schedules
// @Accept  json
// @Produce  json
// @Success 200 {object} util.JobScheduleResult
// @Header 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /listSchedules [get]
func ListSchedules(w http.ResponseWriter, r *http.Request) {
	var jobScheduleResult util.JobScheduleResult
	var jobSchedules []util.JobSchedule
	var result util.Result
	var messages []util.Message

	jobScheduleFiles, err := FindJobSchedules()
	if err != nil {
		msg := util.SetMessage("ERROR", "Couldn't find job schedules! "+err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		jobScheduleResult.Result = result

		_ = json.NewDecoder(r.Body).Decode(&jobScheduleResult)
		json.NewEncoder(w).Encode(jobScheduleResult)

		return
	}

	for _, path := range jobScheduleFiles {
		schedule, err := ReadJobSchedule(path)
		if err != nil {
			msg := util.SetMessage("ERROR", "Couldn't read schedule! "+err.Error())
			messages = append(messages, msg)

			result.Code = 1
			result.Messages = messages

			jobScheduleResult.Result = result

			_ = json.NewDecoder(r.Body).Decode(&jobScheduleResult)
			json.NewEncoder(w).Encode(jobScheduleResult)

			return
		}

		jobSchedules = append(jobSchedules, schedule)

	}
	result.Code = 0
	jobScheduleResult.Result = result
	jobScheduleResult.JobSchedules = jobSchedules

	_ = json.NewDecoder(r.Body).Decode(&jobScheduleResult)
	json.NewEncoder(w).Encode(jobScheduleResult)
}
