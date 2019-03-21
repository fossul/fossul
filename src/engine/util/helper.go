package util

import (
	"time"
	"fmt"
	"strconv"
	"encoding/gob"
	"os"
	"log"
)


func GetTimestamp() int64 {
	time := time.Now().Unix()

	return time
}

func GetBackupDir(configMap map[string]string) string {
	backupPath := configMap["BackupDestPath"] + "/" + configMap["ProfileName"] + "/" + configMap["ConfigName"]

	return backupPath
}

func GetBackupPath(configMap map[string]string) string {
	backupName := GetBackupName(configMap["BackupName"],configMap["BackupPolicy"])
	backupPath := configMap["BackupDestPath"] + "/" + configMap["ProfileName"] + "/" + configMap["ConfigName"] + "/" + backupName

	return backupPath
}

func GetBackupName(name, policy string) string {
	time := GetTimestamp()
	timeToString := fmt.Sprintf("%d",time)

	backupName := fmt.Sprintf(name + "_" + policy + "_" + timeToString)

	return backupName
}

func ConvertEpoch(epoch string) string {
	i := StringToInt64(epoch)
	time:= time.Unix(i,0)

	return time.String()
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

func CreateDir(path string, mode os.FileMode) {

	if ExistsPath(path) == false {
		if err := os.MkdirAll(path, mode); err != nil {
			log.Println("Creating directory " + path + " failed")
			log.Println(err.Error())
			os.Exit(1)
		 } else {
			log.Println("Creating directory " + path + " completed successfully")
		 }		
	}
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
    	log.Println(err.Error())
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
