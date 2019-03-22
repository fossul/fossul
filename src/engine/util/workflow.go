package util

import (
	"log"
	"math/rand"
	"encoding/json"
	"net/http"
	"time"
)

type Workflow struct {
	Id int `json:"id"`
	Status string `json:"status"`
	Steps []Step `json:"steps,omitempty"`
}

type WorkflowResult struct {
	Id int `json:"id"`
	Result Result `json:"result,omitempty"`
}

//func New() *Workflow {
//	w := &Workflow{}
//	return w
//}

type Step struct {
	Id int `json:"id"`
	Status string `json:"status"`
	Label string `json:"label,omitempty"`
}

func SetWorkflowId() int {
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(10000)
	return id
}

func SetStep(id int,status,label string) Step {
	var step Step
	step.Id = id
	step.Status = status
	step.Label = label

	return step
}

func SetWorkflowStatusStart(workflow *Workflow) *Workflow {
	workflow.Status = "RUNNING"
	return workflow
}

func SetWorkflowStatusEnd(workflow *Workflow) *Workflow {
	workflow.Status = "COMPLETE"
	return workflow
}

func SetWorkflowStatusError(workflow *Workflow) *Workflow {
	workflow.Status = "ERROR"
	return workflow
}

func SerializeWorkflow(resultsDir string,workflow *Workflow) {
	CreateDir(resultsDir,0755)
	err := WriteGob(resultsDir + "/workflow",workflow)
	if err != nil {
		log.Println(err.Error())
	}
}

func SerializeWorkflowStepResults(resultsDir string,stepId int, results Result) {
	stepIdToString := IntToString(stepId)
	CreateDir(resultsDir,0755)
	err := WriteGob(resultsDir + "/" + stepIdToString,results)
	if err != nil {
		log.Println(err.Error())
	}
}


func SetWorkflowStep(workflow *Workflow, step Step) *Workflow {
	steps := workflow.Steps
	steps = append(steps, step)

	workflow.Steps = steps
	return workflow
}

func GetWorkflowSteps(w http.ResponseWriter, r *http.Request) []Step {

	var steps []Step
	if err := json.NewDecoder(r.Body).Decode(&steps); err != nil {
		log.Println(err)
	}
	defer r.Body.Close()
 
	res,err := json.Marshal(&steps)
	if err != nil {
        log.Println(err)
	}

	log.Println("DEBUG", string(res))

	return steps
}

