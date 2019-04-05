package main

import (
	"engine/util"
)

type appPlugin string
var config util.Config
var AppPlugin appPlugin

func (a appPlugin) SetEnv(c util.Config) util.Result {
	config = c
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	result = util.SetResult(resultCode, messages)

	return result
}

func (a appPlugin) Quiesce() util.Result {	
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Quiesce ***")
	messages = append(messages,msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (a appPlugin) Unquiesce() util.Result {	
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Unquiesce ***")
	messages = append(messages,msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (a appPlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "sample"
	plugin.Description = "A sample plugin"
	plugin.Type = "app"

	var capabilities []util.Capability
	var quiesceCap util.Capability
	quiesceCap.Name = "quiesce"

	var unquiesceCap util.Capability
	unquiesceCap.Name = "unquiesce"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities,quiesceCap,unquiesceCap,infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func main() {}