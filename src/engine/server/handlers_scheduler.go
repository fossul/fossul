package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"fossil/src/engine/util"
	"log"
)

func AddSchedule(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var policy string = params["policy"]	

	var result util.Result
	var messages []util.Message

	cronSchedule,err := util.GetCronSchedule(w,r)
	if err != nil {
		msg := util.SetMessage("ERROR","Couldn't get cron schedule! " + err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages	
		
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
		
		return
	}

	id,err := AddCronSchedule(profileName,configName,policy,cronSchedule)
	if err != nil {
		msg := util.SetMessage("ERROR","Add schedule failed! " + err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages	
		
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
		
		return
	}

	log.Println("INFO: Added schedule with id", id)

	msg := util.SetMessage("INFO","Add schedule completed successfully")
	messages = append(messages, msg)

	result.Code = 0
	result.Messages = messages	
	
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var profileName string = params["profileName"]
	var configName string = params["configName"]
	var policy string = params["policy"]

	var result util.Result
	var messages []util.Message

	err := DeleteCronSchedule(profileName,configName,policy)
	if err != nil {
		msg := util.SetMessage("ERROR","Delete schedule failed! " + err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages	
		
		_ = json.NewDecoder(r.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
		
		return
	}

	msg := util.SetMessage("INFO","Delete schedule completed successfully")
	messages = append(messages, msg)

	result.Code = 0
	result.Messages = messages	
	
	_ = json.NewDecoder(r.Body).Decode(&result)
	json.NewEncoder(w).Encode(result)
}

func ListSchedules(w http.ResponseWriter, r *http.Request) {
	var jobScheduleResult util.JobScheduleResult
	var jobSchedules []util.JobSchedule
	var result util.Result
	var messages []util.Message

	jobScheduleFiles,err := FindJobSchedules()
	if err != nil {
		msg := util.SetMessage("ERROR","Couldn't find job schedules! " + err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages	

		jobScheduleResult.Result = result
		
		_ = json.NewDecoder(r.Body).Decode(&jobScheduleResult)
		json.NewEncoder(w).Encode(jobScheduleResult)
		
		return
	}

	for _,path := range jobScheduleFiles {
		schedule,err := ReadJobSchedule(path)
		if err != nil {
			msg := util.SetMessage("ERROR","Couldn't read schedule! " + err.Error())
			messages = append(messages, msg)
	
			result.Code = 1
			result.Messages = messages	
	
			jobScheduleResult.Result = result
			
			_ = json.NewDecoder(r.Body).Decode(&jobScheduleResult)
			json.NewEncoder(w).Encode(jobScheduleResult)
			
			return
		}

		jobSchedules = append(jobSchedules,schedule)

	}	
	result.Code = 0
	jobScheduleResult.Result = result
	jobScheduleResult.JobSchedules = jobSchedules

	_ = json.NewDecoder(r.Body).Decode(&jobScheduleResult)
	json.NewEncoder(w).Encode(jobScheduleResult)
}