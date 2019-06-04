package util

import (
	"fmt"
	"os/exec"
	"strings"
)

func ExecuteCommand(args ...string) (result Result) {
	baseCmd := args[0]
	cmdArgs := args[1:]

	var messages []Message
	s0 := fmt.Sprintf("Executing command [%s %s]", baseCmd, strings.Join(cmdArgs, " "))
	message := SetMessage("CMD", s0)
	messages = append(messages, message)

	cmd := exec.Command(baseCmd, cmdArgs...)

	stdoutStderrBytes, err := cmd.CombinedOutput()
	var resultCode int
	if err != nil {
		s1 := fmt.Sprintf("Command [%s %s] failed", baseCmd, strings.Join(cmdArgs, " "))
		message := SetMessage("ERROR", s1)
		messages = append(messages, message)

		s2 := fmt.Sprintf("Command failed with [%s]", err.Error())
		message = SetMessage("ERROR", s2)
		messages = append(messages, message)

		message = SetMessage("ERROR", string(stdoutStderrBytes))
		messages = append(messages, message)

		resultCode = 1
	} else {
		message = SetMessage("INFO", string(stdoutStderrBytes))
		messages = append(messages, message)

		s1 := fmt.Sprintf("Command [%s %s] completed successfully", baseCmd, strings.Join(cmdArgs, " "))
		message := SetMessage("INFO", s1)
		messages = append(messages, message)

		resultCode = 0
	}

	result = SetResult(resultCode, messages)

	return result
}
