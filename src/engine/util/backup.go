package util

type Backups struct {
	Backups        []Backup `json:"backup,omitempty"`
}

type Backup struct {
	Name          string    `json:"name,omitempty"`
	Timestamp     string    `json:"timestamp,omitempty"`
	Epoch         int    	`json:"epoch,omitempty"`
	Policy        string    `json:"policy,omitempty"`
	WorkflowId	  string 	`json:"workflowId,omitempty"`
}

type ByEpoch []Backup

func (a ByEpoch) Len() int           { return len(a) }
func (a ByEpoch) Less(i, j int) bool { return a[i].Epoch < a[j].Epoch }
func (a ByEpoch) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func GetBackupsByPolicy(policy string, backups []Backup) []Backup {
	var backupsByPolicy []Backup
	for _, backup := range backups {
		if policy == backup.Policy {
			backupsByPolicy = append(backupsByPolicy,backup)
		}
	}

	return backupsByPolicy
}