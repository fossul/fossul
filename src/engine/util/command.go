package util

import (
//	"bytes"
	"fmt"
	"os/exec"
	"strings"
//	"reflect"
)

func ExecuteCommand(args ...string) (result Result) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	s := fmt.Sprintf("CMD executing command [%s %s]",baseCmd, strings.Join(cmdArgs, " "))
	fmt.Println(s)

	cmd := exec.Command(baseCmd, cmdArgs...)

	var resultCode int
	stdoutStderrBytes, err := cmd.CombinedOutput()
	if err != nil {
		s1 := fmt.Sprintf("ERROR command [%s %s] failed",baseCmd, strings.Join(cmdArgs, " "))
		fmt.Println(s1)
		s2 := fmt.Sprintf("ERROR command failed with [%s]", err)
		fmt.Println(s2)

		resultCode = 1
	} else {
		resultCode = 0
		s := fmt.Sprintf("INFO command [%s %s] completed successfully",baseCmd, strings.Join(cmdArgs, " "))
		fmt.Println(s)

	}

	output := string(stdoutStderrBytes)
	outputArray := strings.Split(output, "\n")

	messages := SetMessages(outputArray)

	result = SetResult(resultCode, messages)

	return result
}