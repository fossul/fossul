package util

type Result struct {
	Code    int      `json:"code"`
	Stdout  string   `json:"stdout"`
	Stderr  string   `json:"stderr"`
}

type Status struct {
	Msg string `json:"msg"`
}
