package pluginUtil

import (
	"os"
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