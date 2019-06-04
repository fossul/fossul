package util

import (
	"io/ioutil"
	"sort"
	"strings"
)

type Jobs struct {
	Jobs   []Job  `json:"backup,omitempty"`
	Result Result `json:"result,omitempty"`
}

type Job struct {
	Id        int    `json:"id,omitempty"`
	Status    string `json:"status,omitempty"`
	Type      string `json:"type,omitempty"`
	Policy    string `json:"policy"`
	Timestamp string `json:"timestamp,omitempty"`
}

func ListJobs(jobsDir string) ([]Job, error) {
	var jobs []Job
	files, err := ioutil.ReadDir(jobsDir)
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Unix() > files[j].ModTime().Unix()
	})

	for _, f := range files {
		if !strings.Contains(f.Name(), "jobSchedule_") {
			var job Job
			workflowFile := jobsDir + "/" + f.Name() + "/workflow"

			workflow := &Workflow{}
			err := ReadGob(workflowFile, &workflow)
			if err != nil {
				return nil, err
			}

			job.Id = workflow.Id
			job.Status = workflow.Status
			job.Type = workflow.Type
			job.Policy = workflow.Policy
			job.Timestamp = workflow.Timestamp

			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func DeleteJobs(baseDir, profileName, configName string, jobRetention int) Result {
	var result Result
	var messages []Message
	jobsDir := baseDir + "/" + profileName + "/" + configName

	files, err := ioutil.ReadDir(jobsDir)
	if err != nil {
		msg := SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result.Code = 1
		result.Messages = messages

		return result
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Unix() > files[j].ModTime().Unix()
	})

	count := 1
	for _, f := range files {
		if strings.Contains(f.Name(), "jobSchedule_") {
			continue
		}

		jobPath := jobsDir + "/" + f.Name()
		if count > jobRetention {
			msg := SetMessage("INFO", "Job cleanup, more jobs exist for profile ["+profileName+"] config ["+configName+"] than job retention ["+IntToString(jobRetention)+"], deleting job ["+jobPath+"]")
			messages = append(messages, msg)

			err := RecursiveDirDelete(jobPath)
			if err != nil {
				msg := SetMessage("ERROR", err.Error())
				messages = append(messages, msg)

				result.Code = 1
				result.Messages = messages

				return result
			}
		}
		count = count + 1
	}

	result.Code = 0
	result.Messages = messages

	return result
}
