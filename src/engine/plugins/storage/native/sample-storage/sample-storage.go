package main

import (
	"fossil/src/engine/util"
)

type storagePlugin string

var config util.Config
var StoragePlugin storagePlugin

func (s storagePlugin) SetEnv(c util.Config) util.Result {
	config = c
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	result = util.SetResult(resultCode, messages)

	return result
}

func (s storagePlugin) Backup() util.Result {	
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Backup ***")
	messages = append(messages,msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) Restore() util.Result {	
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Restore ***")
	messages = append(messages,msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) BackupDelete() util.Result {	
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Backup Delete ***")
	messages = append(messages,msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) BackupList() util.Backups {	
	var backups util.Backups

	return backups
}

func (s storagePlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "sample-storage"
	plugin.Description = "A sample storage plugin"
	plugin.Version = "1.0.0"
	plugin.Type = "storage"

	var capabilities []util.Capability
	var backupCap util.Capability
	backupCap.Name = "backup"

	var backupListCap util.Capability
	backupListCap.Name = "backupList"

	var backupDeleteCap util.Capability
	backupDeleteCap.Name = "backupDelete"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities,backupCap,backupListCap,backupDeleteCap,infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func main() {}