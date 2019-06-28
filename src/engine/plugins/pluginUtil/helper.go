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
	"fossul/src/engine/util"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
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

	re := regexp.MustCompile(`(\S+)_(\S+)_(\S+)_(\S+)`)
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
	type timeSlice []util.Backup

	re := regexp.MustCompile(`(\S+)_(\S+)_(\S+)_(\S+)`)
	for _, snapshot := range snapshots {
		var backup util.Backup
		match := re.FindStringSubmatch(snapshot)

		if len(match) != 0 {
			if strings.Contains(match[1], backupName) {
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
	}

	sort.Sort(util.ByEpochBackup(backups))

	return backups, nil
}

func ListArchives(dirs []string) ([]util.Archive, error) {

	var archives []util.Archive
	type timeSlice []util.Archive

	re := regexp.MustCompile(`(\S+)_(\S+)_(\S+)_(\S+)`)
	for _, dir := range dirs {
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

			archives = append(archives, archive)
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
