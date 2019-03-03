package util

type Config struct {
	Profile string `json:"profile"`
	AppPlugin string `json:"appPlugin"`
	StoragePlugin string `json:"storagePlugin"`
	BackupRententionDays int `json:"backupRetentionDays"`

}