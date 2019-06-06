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
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Workflow struct {
	Id        int    `json:"id"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Policy    string `json:"policy"`
	Timestamp string `json:"timestamp,omitempty"`
	Steps     []Step `json:"steps,omitempty"`
}

type WorkflowResult struct {
	Id     int    `json:"id"`
	Result Result `json:"result,omitempty"`
}

type WorkflowStatusResult struct {
	Workflow Workflow `json:"workflow,omitempty"`
	Result   Result   `json:"result,omitempty"`
}

type Step struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
	Label  string `json:"label,omitempty"`
}

func GetWorkflowId() int {
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(10000)
	return id
}

func CreateStep(workflow *Workflow) Step {
	id := len(workflow.Steps)

	var step Step
	step.Id = id
	step.Status = "RUNNING"
	step.Label = "Step " + IntToString(id)

	return step
}

func CreateCommentStep(workflow *Workflow) Step {
	id := len(workflow.Steps)

	var step Step
	step.Id = id
	step.Status = "COMPLETE"
	step.Label = "Step " + IntToString(id)

	return step
}

func SetStepComplete(workflow *Workflow, step Step) {
	workflow.Steps[step.Id].Status = "COMPLETE"
}

func SetStepError(workflow *Workflow, step Step) {
	workflow.Steps[step.Id].Status = "ERROR"
}

func SetWorkflowStatusStart(workflow *Workflow) {
	workflow.Status = "RUNNING"
}

func SetWorkflowStatusEnd(workflow *Workflow) {
	workflow.Status = "COMPLETE"
}

func SetWorkflowStatusError(workflow *Workflow) {
	workflow.Status = "ERROR"
}

func SerializeWorkflow(resultsDir string, workflow *Workflow) {
	err := CreateDir(resultsDir, 0755)
	if err != nil {
		log.Println(err.Error())
	}

	err = WriteGob(resultsDir+"/workflow", workflow)
	if err != nil {
		log.Println(err.Error())
	}
}

func SerializeWorkflowStepResults(resultsDir string, stepId int, results Result) {
	stepIdToString := IntToString(stepId)
	err := CreateDir(resultsDir, 0755)
	if err != nil {
		log.Println(err.Error())
	}

	err = WriteGob(resultsDir+"/"+stepIdToString, results)
	if err != nil {
		log.Println(err.Error())
	}
}

func SetWorkflowStep(workflow *Workflow, step Step) {
	steps := workflow.Steps
	steps = append(steps, step)
	workflow.Steps = steps
}

func GetWorkflowSteps(w http.ResponseWriter, r *http.Request) []Step {

	var steps []Step
	if err := json.NewDecoder(r.Body).Decode(&steps); err != nil {
		log.Println(err)
	}
	defer r.Body.Close()

	_, err := json.Marshal(&steps)
	if err != nil {
		log.Println(err)
	}

	//log.Println("DEBUG", string(res))

	return steps
}
