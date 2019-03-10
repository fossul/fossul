package util

import "time"

type Result struct {
	Code   int    `json:"code,omitempty"`
	Messages []string `json:"messages,omitempty"`
	Time   string `json:"time,omitempty"`
}

func SetResult(code int, messages []string) Result {
	t := time.Now()
	
	var result Result
	result.Code = code
	result.Messages = messages
	result.Time = t.Format(time.RFC3339)

	return result
}
