package util

type Result struct {
	Code   int    `json:"code"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}
