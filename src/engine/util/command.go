package util

import (
	"bytes"
	"log"
	"os/exec"
	"os"
	"strings"
//	"reflect"
)

func ExecuteCommand(config Config, pluginType string, args ...string) (result Result) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	log.Println("Executing command: ", baseCmd, cmdArgs)

	cmd := exec.Command(baseCmd, cmdArgs...)

	if pluginType == "app" {
		cmd = setAppPluginEnv(config, cmd)
	} else if pluginType == "storage" {
		cmd = setStoragePluginEnv(config, cmd)
	}

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

func setAppPluginEnv(config Config, cmd *exec.Cmd) *exec.Cmd {
	cmd.Env = os.Environ()
	for k, v := range config.AppPluginParameters { 
		cmd.Env = append(cmd.Env, k + "=" + v)
		log.Println("test", k, v)
	}
	
	return cmd
}

func setStoragePluginEnv(config Config, cmd *exec.Cmd) *exec.Cmd {
	cmd.Env = os.Environ()
	for k, v := range config.StoragePluginParameters { 
		cmd.Env = append(cmd.Env, k + "=" + v)
		log.Println("test", k, v)
	}
	
	return cmd
}