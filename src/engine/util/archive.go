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

type Archives struct {
	Archives []Archive `json:"archive,omitempty"`
	Result   Result    `json:"result,omitempty"`
}

type Archive struct {
	Name       string `json:"name,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
	Epoch      int    `json:"epoch,omitempty"`
	Policy     string `json:"policy,omitempty"`
	WorkflowId string `json:"workflowId,omitempty"`
}
