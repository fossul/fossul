package main

import (
	"bytes"
	"log"
	"os/exec"
    "engine/util"

)

func executeCommand(args ...string) (result util.Result) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	log.Println("Executing command: ", baseCmd, cmdArgs)

	cmd := exec.Command(baseCmd, cmdArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Println("cmd.Run() failed with\n", err)
		result.Code = 1
	} else {
		result.Code = 0
	}

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	result.Stderr = errStr
	result.Stdout = outStr

	return result
}
