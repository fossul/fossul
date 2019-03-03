package util

type Config struct {
	Profile string `json:"profile"`
	AppPlugin string `json:"appPlugin"`
	StoragePlugin string `json:"storagePlugin"`
	BackupRetentions []BackupRetention `json:"backupRetentions"`

}

type BackupRetention struct {
	Policy string `json:"policy"`
	RetentionDays int `json:"retentionDays"`	
}