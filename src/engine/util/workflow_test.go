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
	"log"
	"regexp"
	"testing"
)

func TestWorkflow(t *testing.T) {
	workflow := &Workflow{}
	workflow.Id = GetWorkflowId()
	SetWorkflowStatusStart(workflow)

	step1 := CreateStep(workflow)
	SetWorkflowStep(workflow, step1)

	step2 := CreateStep(workflow)
	SetWorkflowStep(workflow, step2)

	step3 := CreateStep(workflow)
	SetWorkflowStep(workflow, step3)

	SetStepComplete(workflow, step2)
	SetStepError(workflow, step3)

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
