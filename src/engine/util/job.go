package util

import (
	"io/ioutil"
	"sort"
	"strings"
)

type Jobs struct {
	Jobs []Job `json:"backup,omitempty"`
	Result Result `json:"result,omitempty"`
}

type Job struct {
	Id int `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
	Timestamp string    `json:"timestamp,omitempty"`
}

func ListJobs(jobsDir string) ([]Job,error) {
	var jobs []Job
	files, err := ioutil.ReadDir(jobsDir)
	if err != nil {
		return nil,err
	}

	sort.Slice(files, func(i,j int) bool{
		return files[i].ModTime().Unix() > files[j].ModTime().Unix()
	})

	for _, f := range files {
		if ! strings.Contains(f.Name(), "jobSchedule_") {
			var job Job
			workflowFile := jobsDir + "/" + f.Name() + "/workflow"

			workflow := &Workflow{}
			err := ReadGob(workflowFile,&workflow)
			if err != nil {
				return nil,err
			}

			job.Id = workflow.Id
			job.Status = workflow.Status
			job.Timestamp = workflow.Timestamp

			jobs = append(jobs,job)
		}	
	}
	
	return jobs,nil
}