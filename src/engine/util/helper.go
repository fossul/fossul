package util

import (
	"time"
	"fmt"
	"strconv"
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
	backupName := GetBackupName(configMap["BackupName"])
	backupPath := configMap["BackupDestPath"] + "/" + configMap["ProfileName"] + "/" + configMap["ConfigName"] + "/" + backupName

	return backupPath
}

func GetBackupName(name string) string {
	time := GetTimestamp()
	timeToString := fmt.Sprintf("%d",time)

	backupName := fmt.Sprintf(name + "_" + timeToString)

	return backupName
}

func ConvertEpoch(epoch string) string {
	i, err := strconv.ParseInt(epoch, 10, 64)
	if err != nil {
    	fmt.Println(err.Error())
	}

	time:= time.Unix(i,0)

	return time.String()
}