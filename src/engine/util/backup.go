package util

type Backups struct {
	Backups        []Backup `json:"backup,omitempty"`
}

type Backup struct {
	Name          string    `json:"name,omitempty"`
	Timestamp     string    `json:"timestamp,omitempty"`
	Epoch         string    `json:"epoch,omitempty"`
}