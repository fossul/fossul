package util

// NOT CURRENTLY USED, for testing

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"reflect"
)

func ExecutePlugin(config Config, pluginType string, args ...string) (result Result) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	log.Println("Executing command: ", baseCmd, cmdArgs)

	cmd := exec.Command(baseCmd, cmdArgs...)

	if pluginType == "app" {
		cmd = setAppPluginEnv(config, cmd)
	} else if pluginType == "storage" {
		cmd = setStoragePluginEnv(config, cmd)
	}

	stdoutStderrBytes, err := cmd.CombinedOutput()
	var resultCode int
	if err != nil {
		log.Println("cmd.Run() failed with\n", err)
		resultCode = 1
	} else {
		resultCode = 0
	}
	output := string(stdoutStderrBytes)
	messages := strings.Split(output, "\n")
	for index,line := range messages{
		log.Println("test12345", index, line)
	}	

	//var stdout, stderr bytes.Buffer
	//cmd.Stdout = &stdout
	//cmd.Stderr = &stderr

	//var resultCode int
	//err := cmd.Run()
	//if err != nil {
	//	log.Println("cmd.Run() failed with\n", err)
	//	resultCode = 1
	//} else {
	//	resultCode = 0
	//}

	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	//outStr = strings.TrimSuffix(outStr, "\n")
	//errStr = strings.TrimSuffix(errStr, "\n")

	result = SetResult(resultCode, messages)

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

func setBaseContainerPluginEnvtest(config Config) {
	v := reflect.ValueOf(config.BaseContainerPlugin)
	values := make([]interface{}, v.NumField())

    for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
		os.Setenv(v.Type().Field(i).Name,v.Field(i).Interface().(string))
		log.Println("Parsing plugin struct",reflect.TypeOf(v.Field(i).Interface()),v.Type().Field(i).Name,v.Field(i).Interface())
    }
}
