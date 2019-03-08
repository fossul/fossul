package util

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"reflect"
)

func ExecutePlugin(args ...string) (result Result) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	log.Println("Executing command: ", baseCmd, cmdArgs)

	cmd := exec.Command(baseCmd, cmdArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	var resultCode int
	err := cmd.Run()
	if err != nil {
		log.Println("cmd.Run() failed with\n", err)
		resultCode = 1
	} else {
		resultCode = 0
	}

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	outStr = strings.TrimSuffix(outStr, "\n")
	errStr = strings.TrimSuffix(errStr, "\n")

	result = SetResult(resultCode, outStr, errStr)

	return result
}

func setBaseContainerPluginEnv(config Config) {
	v := reflect.ValueOf(config.BaseContainerPlugin)
	values := make([]interface{}, v.NumField())

    for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
		os.Setenv(v.Type().Field(i).Name,v.Field(i).Interface().(string))
		log.Println("Parsing plugin struct",reflect.TypeOf(v.Field(i).Interface()),v.Type().Field(i).Name,v.Field(i).Interface())
    }
}
