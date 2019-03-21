package pluginUtil

import (
	"os"
	"io/ioutil"
	"regexp"
	"engine/util"
	"path/filepath"
	"sort"
)

func ExistsPath(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func CreateDir(path string, mode os.FileMode) {

	if ExistsPath(path) == false {
		if err := os.MkdirAll(path, mode); err != nil {
			LogErrorMessage("creating directory " + path + " failed")
			LogErrorMessage(err.Error())
			os.Exit(1)
		 } else {
			LogInfoMessage("creating directory " + path + " completed successfully")
		 }		
	}
}

func ListBackups(path string) []util.Backup {
	files, err := ioutil.ReadDir(path)
    if err != nil {
		LogErrorMessage(err.Error())
		os.Exit(1)
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

	return backups
}

func RecursiveDirDelete(dir string) {
	if ExistsPath(dir) == true {
		d, err := os.Open(dir)

		if err != nil {
			LogErrorMessage(err.Error())
			os.Exit(1)
		}
		defer d.Close()

		names, err := d.Readdirnames(-1)
		if err != nil {
			LogErrorMessage(err.Error())
			os.Exit(1)
		}

		for _, name := range names {
			err = os.RemoveAll(filepath.Join(dir, name))
			if err != nil {
				LogErrorMessage(err.Error())
				os.Exit(1)
			}
		}

		err = os.Remove(dir)
		if err != nil {
			LogErrorMessage(err.Error())
			os.Exit(1)
		}
		LogInfoMessage("Removed directory " + dir + " completed successfully")
	}	
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