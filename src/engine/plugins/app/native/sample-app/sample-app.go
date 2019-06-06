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
	"strings"
)

type appPlugin string

var AppPlugin appPlugin

func (a appPlugin) SetEnv(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	result = util.SetResult(resultCode, messages)

	return result
}

func (a appPlugin) Discover(config util.Config) util.DiscoverResult {
	var discoverResult util.DiscoverResult = setDiscoverResult()
	return discoverResult
}

func (a appPlugin) Quiesce(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Application Quiesce ***")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (a appPlugin) Unquiesce(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Application Unquiesce ***")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (a appPlugin) PreRestore(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Application PreRestore ***")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (a appPlugin) PostRestore(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Application PostRestore ***")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (a appPlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setDiscoverResult() (discoverResult util.DiscoverResult) {
	var data []string
	data = append(data, "/path/to/data/file1")
	data = append(data, "/path/to/data/file2")

	var logs []string
	logs = append(logs, "/path/to/logs/file1")
	logs = append(logs, "/path/to/logs/file2")

	var discoverInst1 util.Discover
	discoverInst1.Instance = "inst1"
	discoverInst1.DataFilePaths = data
	discoverInst1.LogFilePaths = logs

	var discoverInst2 util.Discover
	discoverInst2.Instance = "inst2"
	discoverInst2.DataFilePaths = data
	discoverInst2.LogFilePaths = logs

	var discoverList []util.Discover
	discoverList = append(discoverList, discoverInst1)
	discoverList = append(discoverList, discoverInst2)

	var messages []util.Message
	msg := util.SetMessage("INFO", "*** Application Discovery ***")
	messages = append(messages, msg)

	for _, discover := range discoverList {
		dataFiles := strings.Join(discover.DataFilePaths, " ")
		logFiles := strings.Join(discover.LogFilePaths, " ")
		msg := util.SetMessage("INFO", "Instance ["+discover.Instance+"] data files: ["+dataFiles+"] log files: ["+logFiles+"]")
		messages = append(messages, msg)
	}

	result := util.SetResult(0, messages)
	discoverResult.Result = result
	discoverResult.DiscoverList = discoverList

	return discoverResult
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "sample"
	plugin.Description = "A sample plugin"
	plugin.Version = "1.0.0"
	plugin.Type = "app"

	var capabilities []util.Capability
	var discoverCap util.Capability
	discoverCap.Name = "discover"

	var quiesceCap util.Capability
	quiesceCap.Name = "quiesce"

	var unquiesceCap util.Capability
	unquiesceCap.Name = "unquiesce"

	var preRestoreCap util.Capability
	preRestoreCap.Name = "preRestore"

	var postRestoreCap util.Capability
	postRestoreCap.Name = "postRestore"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities, discoverCap, quiesceCap, unquiesceCap, preRestoreCap, postRestoreCap, infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func main() {}
