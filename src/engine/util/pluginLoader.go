package util

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"reflect"
	"fmt"
)

func ExecutePlugin(config Config, pluginType string, args ...string) (result Result) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	s := fmt.Sprintf("CMD executing plugin [%s %s]",baseCmd, strings.Join(cmdArgs, " "))
	fmt.Println(s)

	cmd := exec.Command(baseCmd, cmdArgs...)

	if pluginType == "app" {
		cmd = setAppPluginEnv(config, cmd)
	} else if pluginType == "storage" {
		cmd = setStoragePluginEnv(config, cmd)
	}

	var resultCode int
	stdoutStderrBytes, err := cmd.CombinedOutput()
	if err != nil {
		s := fmt.Sprintf("ERROR plugin command [%s %s] failed",baseCmd, strings.Join(cmdArgs, " "))
		fmt.Println(s)
		fmt.Println("ERROR plugin command failed with\n", err)
		resultCode = 1
	} else {
		resultCode = 0
		s := fmt.Sprintf("INFO plugin command [%s %s] completed successfully",baseCmd, strings.Join(cmdArgs, " "))
		fmt.Println(s)

	}

	output := string(stdoutStderrBytes)
	outputArray := strings.Split(output, "\n")

	messages := SetMessages(outputArray)

	result = SetResult(resultCode, messages)

	return result
}

func ExecutePluginInfo(config Config, pluginType string, args ...string) (result ResultSimple) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	s := fmt.Sprintf("CMD executing plugin [%s %s]",baseCmd, strings.Join(cmdArgs, " "))
	fmt.Println(s)

	cmd := exec.Command(baseCmd, cmdArgs...)

	if pluginType == "app" {
		cmd = setAppPluginEnv(config, cmd)
	} else if pluginType == "storage" {
		cmd = setStoragePluginEnv(config, cmd)
	}

	var resultCode int
	stdoutStderrBytes, err := cmd.CombinedOutput()
	if err != nil {
		s := fmt.Sprintf("ERROR plugin command [%s %s] failed",baseCmd, strings.Join(cmdArgs, " "))
		fmt.Println(s)
		fmt.Println("ERROR plugin command failed with\n", err)
		resultCode = 1
	} else {
		resultCode = 0
		s := fmt.Sprintf("INFO plugin command [%s %s] completed successfully",baseCmd, strings.Join(cmdArgs, " "))
		fmt.Println(s)

	}
	
	output := string(stdoutStderrBytes)
	outputArray := strings.Split(output, "\n")

	result.Code = resultCode
	result.Messages = outputArray

	return result
}

func setAppPluginEnv(config Config, cmd *exec.Cmd) *exec.Cmd {
	cmd.Env = os.Environ()
	for k, v := range config.AppPluginParameters { 
		cmd.Env = append(cmd.Env, k + "=" + v)
	}
	
	return cmd
}

func setStoragePluginEnv(config Config, cmd *exec.Cmd) *exec.Cmd {
	cmd.Env = os.Environ()
	for k, v := range config.StoragePluginParameters { 
		cmd.Env = append(cmd.Env, k + "=" + v)
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
