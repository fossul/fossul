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
	"testing"
)

func TestGetArchiveBackupsByPolicy(t *testing.T) {
	var archive1 Archive
	archive1.Epoch = 2147483647
	archive1.Name = "test"
	archive1.Policy = "daily"
	archive1.Timestamp = "Tue 18 Jun 2019 04:03:16 PM CEST"
	archive1.WorkflowId = "1"

	var archives []Archive
	archives = append(archives, archive1)

	archivesByPolicy := GetArchivesByPolicy("daily", archives)
	if len(archivesByPolicy) == 0 {
		t.Fail()
	}
}
