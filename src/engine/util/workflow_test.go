package util

import (
	"testing"
	"regexp"
	"log"
)

func TestWorkflow(t *testing.T) {
	workflow := &Workflow{}
	workflow.Id =  GetWorkflowId()
	SetWorkflowStatusStart(workflow)

	step1 := CreateStep(workflow)
	SetWorkflowStep(workflow,step1)

	step2 := CreateStep(workflow)
	SetWorkflowStep(workflow,step2)

	step3 := CreateStep(workflow)
	SetWorkflowStep(workflow,step3)


	SetStepComplete(workflow,step2)
	SetStepError(workflow,step3)

	log.Println(workflow)

	re := regexp.MustCompile(`\d+`)
	match := re.FindStringSubmatch(IntToString(workflow.Id))

	if len(match) == 0 {
		t.Fail()
	}

	if len(workflow.Steps) != 3 {
		t.Fail()
	}

	if workflow.Steps[0].Status != "RUNNING" {
		t.Fail()
	}

	if workflow.Steps[1].Status != "COMPLETE" {
		t.Fail()
	}

	if workflow.Steps[2].Status != "ERROR" {
		t.Fail()
	}
}