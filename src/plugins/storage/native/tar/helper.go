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
package main

import (
	"errors"
	"io/ioutil"
	"regexp"
	"sort"

	"github.com/fossul/fossul/src/engine/util"
)

func ListBackups(path string) ([]util.Backup, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var backups []util.Backup
	type timeSlice []util.Backup

	re := regexp.MustCompile(`(\S+)-(\S+)-(\d+)-(\d+).tar.gz`)
	for _, f := range files {
		var backup util.Backup
		match := re.FindStringSubmatch(f.Name())

		if len(match) != 0 {
			backup.Name = match[1]
			backup.Policy = match[2]
			backup.WorkflowId = match[3]

			epoch := util.StringToInt(match[4])
			backup.Epoch = epoch

			timestamp := util.ConvertEpoch(match[4])
			backup.Timestamp = timestamp

			backups = append(backups, backup)
		}
	}

	sort.Sort(util.ByEpochBackup(backups))

	return backups, nil
}

func GetBackupNameFromWorkflowId(workflowId, path string) (string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`(\S+)-(\S+)-(\d+)-(\d+).tar.gz`)
	for _, f := range files {
		match := re.FindStringSubmatch(f.Name())

		if len(match) != 0 {
			if match[3] == workflowId {
				return match[1] + "-" + match[2] + "-" + match[3] + "-" + match[4] + ".tar.gz", nil
			}
		}
	}

	return "", errors.New("ERROR: No backup exists for workflow id [" + workflowId + "]")
}

func GetTarBackupCommand(namespace, podName, tarFile, kubeCmd string, backupSrcFilePaths []string) ([]util.Content, []string) {

	var content util.Content
	var contents []util.Content
	var args []string

	args = append(args, kubeCmd)
	args = append(args, "exec")
	args = append(args, "-n")
	args = append(args, namespace)
	args = append(args, podName)
	args = append(args, "--")
	args = append(args, "tar")
	args = append(args, "czf")
	//args = append(args, "-")
	args = append(args, "/tmp/"+tarFile)
	for _, backupSrcFilePath := range backupSrcFilePaths {
		content.Type = "filesystem"
		content.Source = backupSrcFilePath
		content.Data = tarFile
		contents = append(contents, content)

		args = append(args, backupSrcFilePath)
	}
	//args = append(args, ">")

	return contents, args
}

func GetBackupCopyCommand(namespace, podName, kubeCmd, tarFile, destPath string) []string {

	var args []string

	args = append(args, kubeCmd)
	args = append(args, "cp")
	args = append(args, "-n")
	args = append(args, namespace)
	args = append(args, podName+":"+"/tmp/"+tarFile)
	args = append(args, destPath)

	return args
}

func GetRestoreCopyCommand(namespace, podName, kubeCmd, tarFile string) []string {

	var args []string

	args = append(args, kubeCmd)
	args = append(args, "cp")
	args = append(args, tarFile)
	args = append(args, "-n")
	args = append(args, namespace)
	args = append(args, podName+":"+"/tmp")

	return args
}

func GetRestoreTarCommand(namespace, podName, kubeCmd, tarFile, restoreDir string) []string {

	var args []string

	args = append(args, kubeCmd)
	args = append(args, "exec")
	args = append(args, "-n")
	args = append(args, namespace)
	args = append(args, podName)
	args = append(args, "--")
	args = append(args, "tar")
	args = append(args, "xf")
	args = append(args, tarFile)
	args = append(args, "-C")
	args = append(args, restoreDir)

	return args
}

func GetTarCleanupCommand(namespace, podName, kubeCmd, tarFile string) []string {

	var args []string

	args = append(args, kubeCmd)
	args = append(args, "exec")
	args = append(args, "-n")
	args = append(args, namespace)
	args = append(args, podName)
	args = append(args, "--")
	args = append(args, "rm")
	args = append(args, "/tmp/"+tarFile)

	return args
}

func GetKubeCmd(kubeCmd string) string {

	if kubeCmd == "" {
		kubeCmd = "kubectl"
	}

	return kubeCmd
}
