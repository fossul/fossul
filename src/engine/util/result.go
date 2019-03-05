package util

import "time"

type Result struct {
	Code   int    `json:"code,omitempty"`
	Stdout string `json:"stdout,omitempty"`
	Stderr string `json:"stderr,omitempty"`
	Time   string `json:"time,omitempty"`
}

func SetResult(code int, stdout string, stderr string) Result {
	t := time.Now()
	
	var result Result
	result.Code = code
	result.Stdout = stdout
	result.Stderr = stderr
	result.Time = t.Format(time.RFC3339)

	return result
}
