package util

type DiscoverResult struct {
	DiscoverList []Discover `json:"discoverList,omitempty"`
	Result       Result `json:"result,omitempty"`
}

type Discover struct {
	Instance	string `json:"instance,omitempty"`
	DataFiles	[]string `json:"data,omitempty"`
	LogFiles	[]string `json:"logs,omitempty"`
}