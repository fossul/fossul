package util

import (
	"encoding/json"
	"net/http"
	"log"
	"gopkg.in/robfig/cron.v3"
)

type JobScheduleResult struct {
	JobSchedules []JobSchedule `json:"jobSchedules,omitempty"`
	Result Result `json:"result,omitempty"`
}

type JobSchedule struct {
	CronId cron.EntryID `json:"cronId"`
	CronSchedule string `json:"cronSchedule"`
	ProfileName string `json:"profileName"`
	ConfigName string `json:"configName"`
	BackupPolicy string `json:"backupPolicy"`
}

func GetCronSchedule(w http.ResponseWriter, r *http.Request) (string,error) {

	var cronSchedule string
	if err := json.NewDecoder(r.Body).Decode(&cronSchedule); err != nil {
		log.Println(err)
		return cronSchedule,err
	}
	defer r.Body.Close()
 
	res,err := json.Marshal(&cronSchedule)
	if err != nil {
		log.Println(err)
		return cronSchedule,err
	}

	log.Println("DEBUG", string(res))

	return cronSchedule,nil
}

func GetCronScheduleId(w http.ResponseWriter, r *http.Request) (cron.EntryID,error) {

	var cronScheduleId cron.EntryID
	if err := json.NewDecoder(r.Body).Decode(&cronScheduleId); err != nil {
		log.Println(err)
		return cronScheduleId,err
	}
	defer r.Body.Close()
 
	res,err := json.Marshal(&cronScheduleId)
	if err != nil {
		log.Println(err)
		return cronScheduleId,err
	}

	log.Println("DEBUG", string(res))

	return cronScheduleId,nil
}