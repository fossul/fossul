/*
Copyright 2019 The Fossul Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/fossul/fossul/src/engine/util"
)

type storagePlugin string

var StoragePlugin storagePlugin

func (s storagePlugin) SetEnv(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	result = util.SetResult(resultCode, messages)

	return result
}

func (s storagePlugin) Backup(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Backup ***")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) Restore(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Restore ***")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) BackupDelete(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Backup Delete ***")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) BackupDeleteWorkflow(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Backup Delete Workflow ***")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) BackupList(config util.Config) util.Backups {
	var backups util.Backups

	return backups
}

func (s storagePlugin) Mount(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "Mount not implemented or required")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) Unmount(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "Unmount not implemented or required")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
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
	backupDeleteCap.Name = "backupDeleteWorkflow"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities, backupCap, backupListCap, backupDeleteCap, infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func main() {}
