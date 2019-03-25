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

//another way to instantiate Workflow
//func New() *Workflow {
//	w := &Workflow{}
//	return w
//}

type Step struct {
	Id int `json:"id"`
	Status string `json:"status"`
	Label string `json:"label,omitempty"`
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

func UpdateStep(status string, step Step) {
	step.Status = status
}

func SetStepComplete(workflow *Workflow,step Step) {
	workflow.Steps[step.Id].Status = "COMPLETE"
}

func SetStepError(workflow *Workflow,step Step) {
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
 
	res,err := json.Marshal(&steps)
	if err != nil {
        log.Println(err)
	}

	log.Println("DEBUG", string(res))

	return steps
}

