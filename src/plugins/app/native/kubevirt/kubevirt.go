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
	"github.com/fossul/fossul/src/client/k8s"
	"github.com/fossul/fossul/src/engine/util"
)

type appPlugin string

var AppPlugin appPlugin

func (a appPlugin) SetEnv(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message

	result = util.SetResult(0, messages)
	return result
}

func (a appPlugin) Discover(config util.Config) util.DiscoverResult {
	var discoverResult util.DiscoverResult
	var discoverList []util.Discover
	var discover util.Discover
	var result util.Result
	var messages []util.Message

	vmList, err := k8s.ListVirtualMachines(config.AppPluginParameters["Namespace"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", "Listing virtual machines for namespace ["+config.AppPluginParameters["Namespace"]+"] failed! "+err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)

		discoverResult.Result = result
		return discoverResult
	} else {
		msg := util.SetMessage("INFO", "Listing virtual machines for namespace ["+config.AppPluginParameters["Namespace"]+"] successful")
		messages = append(messages, msg)
	}

	for _, vm := range vmList.Items {
		discover.Instance = vm.Name
		discoverList = append(discoverList, discover)
	}

	result = util.SetResult(0, messages)
	discoverResult.Result = result
	discoverResult.DiscoverList = discoverList

	return discoverResult
}

func (a appPlugin) Quiesce(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "Pausing virtual machine ["+config.AppPluginParameters["VmName"]+"]")
	messages = append(messages, msg)

	err := k8s.PauseVirtualMachine(config.AppPluginParameters["Namespace"], config.AccessWithinCluster, config.AppPluginParameters["VmName"], 60000000000)
	if err != nil {
		msg = util.SetMessage("ERROR", "Pausing virtual machine ["+config.AppPluginParameters["VmName"]+"] failed! "+err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)

		return result
	} else {
		msg = util.SetMessage("INFO", "Pausing virtual machine ["+config.AppPluginParameters["VmName"]+"] successful")
		messages = append(messages, msg)
	}

	result = util.SetResult(resultCode, messages)
	return result
}

func (a appPlugin) Unquiesce(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "Un-Pausing virtual machine ["+config.AppPluginParameters["VmName"]+"]")
	messages = append(messages, msg)

	err := k8s.UnPauseVirtualMachine(config.AppPluginParameters["Namespace"], config.AccessWithinCluster, config.AppPluginParameters["VmName"], 60000000000)
	if err != nil {
		msg = util.SetMessage("ERROR", "Un-Pausing virtual machine ["+config.AppPluginParameters["VmName"]+"] failed! "+err.Error())
		messages = append(messages, msg)
		result = util.SetResult(1, messages)

		return result
	} else {
		msg = util.SetMessage("INFO", "Un-Pausing virtual machine ["+config.AppPluginParameters["VmName"]+"] successful")
		messages = append(messages, msg)
	}

	result = util.SetResult(resultCode, messages)
	return result
}

func (a appPlugin) PreRestore(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message

	result = util.SetResult(0, messages)
	return result
}

func (a appPlugin) PostRestore(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message

	result = util.SetResult(0, messages)
	return result
}

func (a appPlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "kubevirt"
	plugin.Description = "Kubevirt plugin for backing up of virtual machines on kubernetes"
	plugin.Version = "1.0.0"
	plugin.Type = "app"

	var capabilities []util.Capability
	var discoverCap util.Capability
	discoverCap.Name = "discover"

	var quiesceCap util.Capability
	quiesceCap.Name = "quiesce"

	var unquiesceCap util.Capability
	unquiesceCap.Name = "unquiesce"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities, discoverCap, quiesceCap, unquiesceCap, infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func main() {}
