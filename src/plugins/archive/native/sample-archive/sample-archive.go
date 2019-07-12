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
	"fossul/src/engine/util"
)

type archivePlugin string

var ArchivePlugin archivePlugin

func (r archivePlugin) SetEnv(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	result = util.SetResult(resultCode, messages)

	return result
}

func (r archivePlugin) Archive(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Archive ***")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (r archivePlugin) ArchiveDelete(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Archive Delete ***")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (r archivePlugin) ArchiveList(config util.Config) util.Archives {
	var archives util.Archives

	return archives
}

func (r archivePlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "sample-archive"
	plugin.Description = "A sample archive plugin"
	plugin.Version = "1.0.0"
	plugin.Type = "archive"

	var capabilities []util.Capability
	var archiveCap util.Capability
	archiveCap.Name = "archive"

	var archiveListCap util.Capability
	archiveListCap.Name = "archiveList"

	var archiveDeleteCap util.Capability
	archiveDeleteCap.Name = "archiveDelete"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities, archiveCap, archiveListCap, archiveDeleteCap, infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func main() {}
