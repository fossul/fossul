package util

type Archives struct {
	Archives []Archive `json:"archive,omitempty"`
	Result   Result    `json:"result,omitempty"`
}

type Archive struct {
	Name       string `json:"name,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
	Epoch      int    `json:"epoch,omitempty"`
	Policy     string `json:"policy,omitempty"`
	WorkflowId string `json:"workflowId,omitempty"`
}
