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
	"github.com/fossul/fossul/src/client"
	"github.com/fossul/fossul/src/engine/util"
	"gopkg.in/robfig/cron.v3"
	"os"
	"strings"
)

var c *CronScheduler

type CronScheduler struct {
	cronScheduler *cron.Cron
}

func StartCron() {
	c := getCron()
	c.cronScheduler.Start()
}

func AddCronSchedule(profileName, configName, policy, cronSchedule string) (cron.EntryID, error) {
	var id cron.EntryID
	var err error

	auth := SetAuth()

	id, err = c.cronScheduler.AddFunc(cronSchedule, func() {
		client.StartBackupWorkflow(auth, profileName, configName, policy)
	})

	if err != nil {
		return id, err
	}

	err = writeJobSchedule(id, cronSchedule, profileName, configName, policy)
	if err != nil {
		return id, err
	}

	return id, nil
}

func DeleteCronSchedule(profileName, configName, policy string) error {
	path := dataDir + "/" + profileName + "/" + configName + "/jobSchedule_" + policy
	schedule, err := ReadJobSchedule(path)
	if err != nil {
		return err
	}
	c.cronScheduler.Remove(schedule.CronId)

	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func ListCronSchedule() []cron.Entry {
	entries := c.cronScheduler.Entries()
	return entries
}

func LoadCronSchedules() error {
	jobScheduleFiles, err := FindJobSchedules()
	if err != nil {
		return err
	}

	for _, path := range jobScheduleFiles {
		schedule, err := ReadJobSchedule(path)
		if err != nil {
			return err
		}

		id, err := AddCronSchedule(schedule.ProfileName, schedule.ConfigName, schedule.BackupPolicy, schedule.CronSchedule)
		if err != nil {
			return err
		}

		err = writeJobSchedule(id, schedule.CronSchedule, schedule.ProfileName, schedule.ConfigName, schedule.BackupPolicy)
		if err != nil {
			return err
		}
	}
	return nil
}

func getCron() *CronScheduler {
	c = new(CronScheduler)
	c.cronScheduler = cron.New()
	return c
}

func writeJobSchedule(id cron.EntryID, cronSchedule, profileName, configName, policy string) error {
	var jobSchedule util.JobSchedule
	jobSchedule.CronId = id
	jobSchedule.CronSchedule = cronSchedule
	jobSchedule.ProfileName = profileName
	jobSchedule.ConfigName = configName
	jobSchedule.BackupPolicy = policy

	scheduleFileDir := dataDir + "/" + profileName + "/" + configName
	scheduleFile := dataDir + "/" + profileName + "/" + configName + "/jobSchedule_" + policy
	err := util.CreateDir(scheduleFileDir, 0755)
	if err != nil {
		return err
	}

	err = util.WriteGob(scheduleFile, jobSchedule)
	if err != nil {
		return err
	}

	return nil
}

func ReadJobSchedule(scheduleFile string) (util.JobSchedule, error) {
	jobSchedule := &util.JobSchedule{}
	err := util.ReadGob(scheduleFile, &jobSchedule)
	if err != nil {
		return *jobSchedule, err
	}

	return *jobSchedule, nil
}

func FindJobSchedules() ([]string, error) {
	var jobScheduleFiles []string

	profiles, err := util.DirectoryList(dataDir)
	if err != nil {
		return jobScheduleFiles, err
	}

	for _, profile := range profiles {
		configs, err := util.DirectoryList(dataDir + "/" + profile)
		if err != nil {
			return jobScheduleFiles, err
		}

		for _, config := range configs {
			files, err := util.FileList(dataDir + "/" + profile + "/" + config)
			if err != nil {
				return jobScheduleFiles, err
			}

			for _, file := range files {
				if strings.Contains(file, "jobSchedule_") {
					path := dataDir + "/" + profile + "/" + config + "/" + file
					jobScheduleFiles = append(jobScheduleFiles, path)
				}
			}
		}
	}

	return jobScheduleFiles, nil
}
