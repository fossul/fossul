package util

import (
	"encoding/json"
	"gopkg.in/robfig/cron.v3"
	"net/http"
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
