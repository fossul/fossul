package util

type Plugin struct {
	Name   int    `json:"name,omitempty"`
	Capabilities []Capability `json:"capabilities"`
}

type Capability struct {
	Name string `json:"policy"`	
}