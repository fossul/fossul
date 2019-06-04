package util

import (
	"log"
	"strings"
	"testing"
)

func TestExecuteCommand(t *testing.T) {
	cmd := "/usr/bin/ls,-a,/tmp"
	args := strings.Split(cmd, ",")

	result := ExecuteCommand(args...)

	log.Println(result)

	if result.Code != 0 {
		t.Fail()
	}
}
