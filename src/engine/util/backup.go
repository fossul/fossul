package util

type Backups struct {
	Backups        []Backup `json:"backup,omitempty"`
}

type Backup struct {
	Name          string    `json:"name,omitempty"`
	Timestamp     string    `json:"timestamp,omitempty"`
	Epoch         string    `json:"epoch,omitempty"`
	Policy        string    `json:"policy,omitempty"`
}

func GetBackupsByPolicy(policy string, backups []Backup) []Backup {
	var backupsByPolicy []Backup
	for _, backup := range backups {
		if policy == backup.Policy {
			backupsByPolicy = append(backupsByPolicy,backup)
		}
	}

	return backupsByPolicy
}