package util

import (
	"time"
	"fmt"
)


func GetTimestamp() int64 {
	time := time.Now().Unix()

	return time
}

func GetBackupName(name string) string {
	time := GetTimestamp()
	timeToString := fmt.Sprintf("%d",time)

	backupName := fmt.Sprintf(name + "_" + timeToString)

	return backupName
}