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
	"fmt"
	"strings"

	"github.com/fossul/fossul/src/client/k8s"
	"github.com/fossul/fossul/src/engine/util"
)

type appPlugin string

var AppPlugin appPlugin

func (a appPlugin) SetEnv(config util.Config) util.Result {
	var result util.Result

	return result
}

func (a appPlugin) Discover(config util.Config) util.DiscoverResult {
	var discoverResult util.DiscoverResult
	var discoverList []util.Discover
	var discover util.Discover
	var result util.Result
	var messages []util.Message

	discover.Instance = config.AppPluginParameters["MongoDb"]

	var dataFilePaths []string
	dumpPath := config.AppPluginParameters["MongoDumpPath"] + "/" + config.WorkflowId
	dataFilePaths = append(dataFilePaths, dumpPath)
	discover.DataFilePaths = dataFilePaths

	msg := util.SetMessage("INFO", "Data Directory is ["+config.AppPluginParameters["MongoDumpPath"]+"]")
	messages = append(messages, msg)

	discoverList = append(discoverList, discover)

	result = util.SetResult(0, messages)
	discoverResult.Result = result
	discoverResult.DiscoverList = discoverList

	return discoverResult
}

func (a appPlugin) Quiesce(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message

	var args []string
	var mkdirArgs []string

	podName, err := k8s.GetPodName(config.AppPluginParameters["Namespace"], config.AppPluginParameters["PodSelector"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	dumpPath := config.AppPluginParameters["MongoDumpPath"] + "/" + config.WorkflowId

	//create tmp directory for storing dump
	mkdirArgs = append(mkdirArgs, "mkdir")
	mkdirArgs = append(mkdirArgs, "-p")
	mkdirArgs = append(mkdirArgs, dumpPath)

	cmdResult := k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AccessWithinCluster, mkdirArgs...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	//execute database dump
	args = append(args, config.AppPluginParameters["MongoDumpCmd"])
	args = append(args, "--host")
	args = append(args, config.AppPluginParameters["MongoHost"])
	args = append(args, "--port")
	args = append(args, config.AppPluginParameters["MongoPort"])
	args = append(args, "--db")
	args = append(args, config.AppPluginParameters["MongoDb"])
	args = append(args, "--username")
	args = append(args, config.AppPluginParameters["MongoUser"])

	if config.AppPluginParameters["MongoPassword"] != "" {
		args = append(args, "--password")
		args = append(args, config.AppPluginParameters["MongoPassword"])
	}

	args = append(args, "--out")
	args = append(args, dumpPath)

	//args = append(args,"--quiet")

	cmdResult = k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AccessWithinCluster, args...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	result = util.SetResult(0, messages)
	return result
}

func (a appPlugin) Unquiesce(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message

	dumpPath := config.AppPluginParameters["MongoDumpPath"] + "/" + config.WorkflowId

	var args []string
	args = append(args, "rm")
	args = append(args, "-rf")
	args = append(args, dumpPath)

	podName, err := k8s.GetPodName(config.AppPluginParameters["Namespace"], config.AppPluginParameters["PodSelector"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	cmdResult := k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AccessWithinCluster, args...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(cmdResult.Messages, messages)
	}

	result = util.SetResult(0, messages)
	return result
}

func (a appPlugin) PreRestore(config util.Config) util.Result {

	var result util.Result
	var messages []util.Message

	podName, err := k8s.GetPodName(config.AppPluginParameters["Namespace"], config.AppPluginParameters["PodSelector"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	//create tmp directory for storing dump
	var mkdirArgs []string
	restoreDir := "/tmp/" + util.IntToString(config.SelectedWorkflowId)
	mkdirArgs = append(mkdirArgs, "mkdir")
	mkdirArgs = append(mkdirArgs, "-p")
	mkdirArgs = append(mkdirArgs, restoreDir)

	cmdResult := k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AccessWithinCluster, mkdirArgs...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	result = util.SetResult(0, messages)
	return result
}

func (a appPlugin) PostRestore(config util.Config) util.Result {

	var result util.Result
	var messages []util.Message

	podName, err := k8s.GetPodName(config.AppPluginParameters["Namespace"], config.AppPluginParameters["PodSelector"], config.AccessWithinCluster)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	var lsDirArgs []string
	lsDirArgs = append(lsDirArgs, "ls")
	lsDirArgs = append(lsDirArgs, "/tmp/"+util.IntToString(config.SelectedWorkflowId))

	cmdResult, restoreDir := k8s.ExecuteCommandWithStdout(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AccessWithinCluster, lsDirArgs...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	restorePath := "/tmp/" + util.IntToString(config.SelectedWorkflowId) + "/" + strings.TrimSpace(restoreDir) + "/" + util.IntToString(config.SelectedWorkflowId) + "/" + config.AppPluginParameters["MongoDb"]

	//execute database restore
	var restoreArgs []string
	restoreArgs = append(restoreArgs, "/bin/sh")
	restoreArgs = append(restoreArgs, "-c")

	if config.AppPluginParameters["MongoPassword"] != "" {
		restoreArgs = append(restoreArgs, config.AppPluginParameters["MongoRestoreCmd"]+" --host "+config.AppPluginParameters["MongoHost"]+
			" --port "+config.AppPluginParameters["MongoPort"]+" --db "+config.AppPluginParameters["MongoDb"]+" --username "+
			config.AppPluginParameters["MongoUser"]+" --password "+config.AppPluginParameters["MongoPassword"]+" "+restorePath)
	} else {
		restoreArgs = append(restoreArgs, config.AppPluginParameters["MongoRestoreCmd"]+" --host "+config.AppPluginParameters["MongoHost"]+
			" --port "+config.AppPluginParameters["MongoPort"]+" --db "+config.AppPluginParameters["MongoDb"]+" --username "+
			config.AppPluginParameters["MongoPassword"]+" "+restorePath)
	}

	cmdResult = k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AccessWithinCluster, restoreArgs...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	var rmDirArgs []string
	restoreTmpDir := "/tmp/" + util.IntToString(config.SelectedWorkflowId)
	rmDirArgs = append(rmDirArgs, "rm")
	rmDirArgs = append(rmDirArgs, "-rf")
	rmDirArgs = append(rmDirArgs, restoreTmpDir)

	cmdResult = k8s.ExecuteCommand(podName, config.AppPluginParameters["ContainerName"], config.AppPluginParameters["Namespace"], config.AccessWithinCluster, rmDirArgs...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages, cmdResult.Messages)
	}

	result = util.SetResult(0, messages)
	return result
}

func (a appPlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "mongo-dump"
	plugin.Description = "Mongo plugin for backing up Mongo database using mongodump utility"
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

func checkErr(err error) {
	fmt.Println("error handling")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {}
