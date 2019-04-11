package util

type DiscoverResult struct {
	DiscoverList []Discover `json:"discoverList,omitempty"`
	Result       Result `json:"result,omitempty"`
}

type Discover struct {
	Instance		string `json:"instance,omitempty"`
	DataFilePaths	[]string `json:"data,omitempty"`
	LogFilePaths	[]string `json:"logs,omitempty"`
}