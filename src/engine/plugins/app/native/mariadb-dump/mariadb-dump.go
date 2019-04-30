package main

import (
	"fossil/src/engine/util"
	"fossil/src/engine/client/k8s"
	"fmt"
)

type appPlugin string

var config util.Config
var AppPlugin appPlugin

func (a appPlugin) SetEnv(c util.Config) util.Result {
	config = c
	var result util.Result

	return result
}	

func (a appPlugin) Discover() util.DiscoverResult {
	var discoverResult util.DiscoverResult
	var discoverList []util.Discover
	var discover util.Discover
	var result util.Result
	var messages []util.Message

	discover.Instance = config.AppPluginParameters["MysqlDb"]

	var dataFilePaths []string
	dumpPath := config.AppPluginParameters["MysqlDumpPath"] + "/" + config.WorkflowId 
	dataFilePaths = append(dataFilePaths,dumpPath)
	discover.DataFilePaths = dataFilePaths

	msg := util.SetMessage("INFO", "Data Directory is [" + dumpPath + "]")
	messages = append(messages,msg)

	discoverList = append(discoverList,discover)

	result = util.SetResult(0, messages)
	discoverResult.Result = result
	discoverResult.DiscoverList = discoverList

	return discoverResult
}	

func (a appPlugin) Quiesce() util.Result {	
	var result util.Result
	var messages []util.Message

	var args []string
	var mkdirArgs []string

	podName,err := k8s.GetPod(config.AppPluginParameters["Namespace"],config.AppPluginParameters["ServiceName"],config.AppPluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages,msg)

		result = util.SetResult(1, messages)
		return result
	}

	dumpPath := config.AppPluginParameters["MysqlDumpPath"] + "/" + config.WorkflowId 

	//create tmp directory for storing dump
	mkdirArgs = append(mkdirArgs,"mkdir")
	mkdirArgs = append(mkdirArgs,"-p")
	mkdirArgs = append(mkdirArgs,dumpPath)
	
	cmdResult := k8s.ExecuteCommand(podName,config.AppPluginParameters["ContainerName"],config.AppPluginParameters["Namespace"],config.AppPluginParameters["AccessWithinCluster"],mkdirArgs...)
	
	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages,cmdResult.Messages)
	}

	//execute database dump
	args = append(args,config.AppPluginParameters["MysqlDumpCmd"])
	args = append(args,"-h")
	args = append(args,config.AppPluginParameters["MysqlHost"])
	args = append(args,"-P")
	args = append(args,config.AppPluginParameters["MysqlPort"])
	args = append(args,"-u")
	args = append(args,config.AppPluginParameters["MysqlUser"])

	if config.AppPluginParameters["MysqlPassword"] != "" {
		args = append(args,"-p" + config.AppPluginParameters["MysqlPassword"])
	} 	

	args = append(args,config.AppPluginParameters["MysqlDb"])

	args = append(args,"-T")
	args = append(args,dumpPath)

	cmdResult = k8s.ExecuteCommand(podName,config.AppPluginParameters["ContainerName"],config.AppPluginParameters["Namespace"],config.AppPluginParameters["AccessWithinCluster"],args...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(messages,cmdResult.Messages)
	}

	result = util.SetResult(0, messages)
	return result
}

func (a appPlugin) Unquiesce() util.Result {	
	var result util.Result
	var messages []util.Message

	dumpPath := config.AppPluginParameters["MysqlDumpPath"] + "/" + config.WorkflowId 

	var args []string
	args = append(args,"rm")
	args = append(args,"-rf")
	args = append(args,dumpPath)

	podName,err := k8s.GetPod(config.AppPluginParameters["Namespace"],config.AppPluginParameters["ServiceName"],config.AppPluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages,msg)

		result = util.SetResult(1, messages)
		return result
	}

	cmdResult := k8s.ExecuteCommand(podName,config.AppPluginParameters["ContainerName"],config.AppPluginParameters["Namespace"],config.AppPluginParameters["AccessWithinCluster"],args...)

	if cmdResult.Code != 0 {
		return cmdResult
	} else {
		messages = util.PrependMessages(cmdResult.Messages,messages)
	}

	result = util.SetResult(0, messages)
	return result
}

func (a appPlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "mariadb-dump"
	plugin.Description = "MariaDB plugin for backing up MySql or MariaDB using mysqldump utility"
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

	capabilities = append(capabilities,discoverCap,quiesceCap,unquiesceCap,infoCap)

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