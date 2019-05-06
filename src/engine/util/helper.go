package util

import (
	"time"
	"fmt"
	"strconv"
	"encoding/gob"
	"os"
	"log"
	"path/filepath"
)


func GetTimestamp() int64 {
	time := time.Now().Unix()

	return time
}

func GetBackupDirFromMap(configMap map[string]string) string {
	backupPath := configMap["BackupDestPath"] + "/" + configMap["ProfileName"] + "/" + configMap["ConfigName"]

	return backupPath
}

func GetBackupDirFromConfig(config Config) string {
	backupPath := config.StoragePluginParameters["BackupDestPath"] + "/" + config.ProfileName + "/" + config.ConfigName

	return backupPath
}

func GetBackupPathFromMap(configMap map[string]string) string {
	backupName := GetBackupName(configMap["BackupName"],configMap["BackupPolicy"],configMap["WorkflowId"])
	backupPath := configMap["BackupDestPath"] + "/" + configMap["ProfileName"] + "/" + configMap["ConfigName"] + "/" + backupName

	return backupPath
}

func GetBackupPathFromConfig(config Config) string {
	backupName := GetBackupName(config.StoragePluginParameters["BackupName"],config.SelectedBackupPolicy,config.WorkflowId)
	backupPath := config.StoragePluginParameters["BackupDestPath"] + "/" + config.ProfileName + "/" + config.ConfigName + "/" + backupName

	return backupPath
}

func GetBackupName(name, policy, workflowId string) string {
	time := GetTimestamp()
	timeToString := fmt.Sprintf("%d",time)

	backupName := fmt.Sprintf(name + "_" + policy + "_" + workflowId + "_" + timeToString)

	return backupName
}

func ConvertEpoch(epoch string) string {
	i := StringToInt64(epoch)
	time:= time.Unix(i,0)

	return time.String()
}

func ConvertEpochToTime(epoch string) time.Time {
	i := StringToInt64(epoch)
	time:= time.Unix(i,0)

	return time
}

func JoinArray(array, combinedArray []string) []string {
	for _, item := range array {
		combinedArray = append (combinedArray,item)
	}

	return combinedArray
}

func ExistsInArray(array []string, str string) bool {
	for _, item := range array {
	   if item == str {
		  return true
	   }
	}
	return false
 }

 func WriteGob(filePath string,object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		   encoder := gob.NewEncoder(file)
		   encoder.Encode(object)
	}
	file.Close()
	return err
}

func ReadGob(filePath string,object interface{}) error {

	file, err := os.Open(filePath)
	if err == nil {
		   decoder := gob.NewDecoder(file)
		   err = decoder.Decode(object)
	}
	file.Close()
	return err
}

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
			return err
		 } else {
			log.Println("Creating directory " + path + " completed successfully")
			return nil
		 }		
	}
	return nil
}

func RecursiveDirDelete(dir string) error {
	if ExistsPath(dir) == true {
		log.Println("Removing directory " + dir)
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

		log.Println("Removing directory " + dir + " completed successfully")
	}	
	return nil
}

func BoolToString(b bool) string {
	s := strconv.FormatBool(b)

	return s
}

func IntToString(i int) string {
	s := strconv.Itoa(i)

	return s
}

func Int64ToString(i int64) string {
	str := strconv.FormatInt(i, 10)
	return str
}

func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
    	log.Println("ERROR " + err.Error())
	}
	return i
}

func StringToInt64(s string) int64 {
	i,err := strconv.ParseInt(s,10,64)
	if err != nil {
    	log.Println(err.Error())
	}

	return i
}

func IntInSlice(i int, list []int) bool {
	for _, v := range list {
		if v == i {
			return true
		}
	}
	return false
}
