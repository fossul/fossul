package util

type PluginInfoResult struct {
	Plugin	   Plugin `json:"plugin,omitempty"`
	Result         Result `json:"result,omitempty"`
}

type Plugin struct {
	Name          string    `json:"name,omitempty"`
	Description   string    `json:"description,omitempty"`
	Type          string    `json:"type,omitempty"`
	Capabilities  []Capability `json:"capabilities"`
}

type Capability struct {
	Name string `json:"name"`	
}