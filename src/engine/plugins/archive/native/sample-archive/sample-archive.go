package main

import (
	"fossil/src/engine/util"
)

type archivePlugin string

var config util.Config
var ArchivePlugin archivePlugin

func (r archivePlugin) SetEnv(c util.Config) util.Result {
	config = c
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	result = util.SetResult(resultCode, messages)

	return result
}

func (r archivePlugin) Archive() util.Result {	
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Archive ***")
	messages = append(messages,msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (r archivePlugin) ArchiveDelete() util.Result {	
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Archive Delete ***")
	messages = append(messages,msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (r archivePlugin) ArchiveList() []util.Archive {	
	var archives []util.Archive

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

	capabilities = append(capabilities,archiveCap,archiveListCap,archiveDeleteCap,infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func main() {}