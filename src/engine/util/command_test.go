package util

import (
	"testing"
	"strings"
)

func TestExecuteCommand(t *testing.T) {
	cmd := "/usr/bin/ls,-a,/tmp"
	args := strings.Split(cmd, ",")

	result := ExecuteCommand(args...)
			
	if result.Code != 0 {
		t.Fail()
	}
}