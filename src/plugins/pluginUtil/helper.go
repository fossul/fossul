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
package pluginUtil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/fossul/fossul/src/engine/util"
)

func ExistsPath(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func CreateDir(path string, mode os.FileMode) error {

	if ExistsPath(path) == false {
		if err := os.MkdirAll(path, mode); err != nil {
			return errors.New("creating directory " + path + " failed!" + err.Error())
		}
	}
	return nil
}

func ListDir(path string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func ListBackups(path string) ([]util.Backup, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var backups []util.Backup
	type timeSlice []util.Backup

	re := regexp.MustCompile(`(\S+)-(\S+)-(\S+)-(\S+)`)
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

func ListSnapshots(snapshots []string, backupName string) ([]util.Backup, error) {

	var backups []util.Backup
	var contents []util.Content
	var content util.Content
	type timeSlice []util.Backup

	re := regexp.MustCompile(`(\S+)-(\S+)-(\d+)-(\d+)-(\S+)`)
	backupsCountMap := make(map[int]util.Backup)
	for _, snapshot := range snapshots {
		var backup util.Backup
		match := re.FindStringSubmatch(snapshot)

		if len(match) != 0 {
			backup.Name = match[1]
			if backup.Name != backupName {
				continue
			}

			backup.Policy = match[2]
			backup.WorkflowId = match[3]

			epoch := util.StringToInt(match[4])
			backup.Epoch = epoch

			timestamp := util.ConvertEpoch(match[4])
			backup.Timestamp = timestamp

			content.Source = match[5]
			content.Data = snapshot
			content.Type = "volume"

			contents = append(contents, content)
			backup.Contents = contents

			backupsCountMap[epoch] = backup
		}

	}

	for k, _ := range backupsCountMap {
		var sortedContents []util.Content
		backup := backupsCountMap[k]
		backupName := backup.Name + "-" + backup.Policy + "-" + backup.WorkflowId + "-" + util.IntToString(backup.Epoch)
		for _, content = range backup.Contents {
			if strings.Contains(content.Data, backupName) {
				fmt.Println("match " + backup.Name)
				sortedContents = append(sortedContents, content)
			}
		}
		backup.Contents = sortedContents
		backups = append(backups, backup)
	}

	sort.Sort(util.ByEpochBackup(backups))

	return backups, nil
}

func ListArchives(dirs []string) ([]util.Archive, error) {

	var archives []util.Archive
	type timeSlice []util.Archive

	re := regexp.MustCompile(`(\S+)-(\S+)-(\S+)-(\S+)`)
	for _, dir := range dirs {
		var isUnique bool = true
		var archive util.Archive
		match := re.FindStringSubmatch(dir)

		if len(match) != 0 {
			archive.Name = match[1]
			archive.Policy = match[2]
			archive.WorkflowId = match[3]

			epoch := util.StringToInt(match[4])
			archive.Epoch = epoch

			timestamp := util.ConvertEpoch(match[4])
			archive.Timestamp = timestamp

			//ensure backup is unique
			for _, a := range archives {
				if a == archive {
					isUnique = false
				}
			}

			if isUnique {
				archives = append(archives, archive)
			}
		}
	}

	sort.Sort(util.ByEpochArchive(archives))

	return archives, nil
}

func GetDirFromPath(path string) string {
	re := regexp.MustCompile(`\S+\/(\S+)$`)
	match := re.FindStringSubmatch(path)

	dir := match[1]

	return dir
}

func RecursiveDirDelete(dir string) error {
	if ExistsPath(dir) == true {
		d, err := os.Open(dir)

		if err != nil {
			return err
		}
		defer d.Close()

		names, err := d.Readdirnames(-1)
		if err != nil {
			return err
		}

		for _, name := range names {
			err = os.RemoveAll(filepath.Join(dir, name))
			if err != nil {
				return err
			}
		}

		err = os.Remove(dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReverseBackupList(backups []util.Backup) chan util.Backup {
	ret := make(chan util.Backup)
	go func() {
		for i, _ := range backups {
			ret <- backups[len(backups)-1-i]
		}
		close(ret)
	}()
	return ret
}

func ReverseArchiveList(archives []util.Archive) chan util.Archive {
	ret := make(chan util.Archive)
	go func() {
		for i, _ := range archives {
			ret <- archives[len(archives)-1-i]
		}
		close(ret)
	}()
	return ret
}

func CephVolumeName(volumeHandle string) string {
	re := regexp.MustCompile(`\d+-\d+-openshift-storage-\d+-(\S+)$`)
	match := re.FindStringSubmatch(volumeHandle)

	csiPath := match[1]
	volumeName := "csi-vol-" + csiPath

	return volumeName
}

func CephSnapshotName(snapshotHandle string) string {
	re := regexp.MustCompile(`\d+-\d+-openshift-storage-\d+-(\S+)$`)
	match := re.FindStringSubmatch(snapshotHandle)

	csiPath := match[1]
	snapshotName := "csi-snap-" + csiPath

	return snapshotName
}
