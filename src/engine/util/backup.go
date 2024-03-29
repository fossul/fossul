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
	"fmt"
	"log"
	"net/http"
)

type Backups struct {
	Backups []Backup `json:"backup,omitempty"`
	Result  Result   `json:"result,omitempty"`
}

type BackupByWorkflow struct {
	Backup Backup `json:"backup,omitempty"`
	Result Result `json:"result,omitempty"`
}

type Backup struct {
	Name       string    `json:"name,omitempty"`
	Timestamp  string    `json:"timestamp,omitempty"`
	Epoch      int       `json:"epoch,omitempty"`
	Policy     string    `json:"policy,omitempty"`
	WorkflowId string    `json:"workflowId,omitempty"`
	Contents   []Content `json:"contents,omitempty"`
}

type Content struct {
	Type                string `json:"type,omitempty"`
	Source              string `json:"source,omitempty"`
	Metadata            string `json:"metadata,omitempty"`
	Data                string `json:"data,omitempty"`
	Size                string `json:"size,omitempty"`
	StorageClass        string `json:"storageClass,omitempty"`
	VolumeSnapshotClass string `json:"volumeSnapshotClass,omitempty"`
}

type ByEpochBackup []Backup

func (a ByEpochBackup) Len() int           { return len(a) }
func (a ByEpochBackup) Less(i, j int) bool { return a[i].Epoch < a[j].Epoch }
func (a ByEpochBackup) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func GetBackupsByPolicy(policy string, backups []Backup) []Backup {
	var backupsByPolicy []Backup
	for _, backup := range backups {
		if policy == backup.Policy {
			backupsByPolicy = append(backupsByPolicy, backup)
		}
	}

	return backupsByPolicy
}

func SerializeBackup(resultsDir string, backup *Backup) {
	err := CreateDir(resultsDir, 0755)
	if err != nil {
		log.Println(err.Error())
	}

	err = WriteGob(resultsDir+"/backup", backup)
	if err != nil {
		log.Println(err.Error())
	}
}

func GetBackup(w http.ResponseWriter, r *http.Request) (Backup, error) {

	var backup Backup
	fmt.Println(r.Body)
	if err := json.NewDecoder(r.Body).Decode(&backup); err != nil {
		return backup, err
	}

	defer r.Body.Close()

	_, err := json.Marshal(&backup)
	if err != nil {
		return backup, err
	}

	return backup, nil
}
