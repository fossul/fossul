package pluginUtil

import (
	"os"
	"io/ioutil"
	"regexp"
	"fossil/src/engine/util"
	"path/filepath"
	"sort"
	"errors"
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

func ListBackups(path string) ([]util.Backup, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil,err
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

	sort.Sort(util.ByEpoch(backups))

	return backups,nil
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
