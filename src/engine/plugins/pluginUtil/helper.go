package pluginUtil

import (
	"os"
	"io/ioutil"
	"regexp"
	"engine/util"
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
	re := regexp.MustCompile(`(\S+)_(\S+)`)
    for _, f := range files {
		var backup util.Backup
		match := re.FindStringSubmatch(f.Name())
		if len(match) != 0 {
			timestamp := util.ConvertEpoch(match[2])
			backup.Name = match[1]
			backup.Timestamp = timestamp
			backups = append(backups, backup)
		}	
	}
	return backups
}