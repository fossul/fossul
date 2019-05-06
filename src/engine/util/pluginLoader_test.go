package util

import (
	"testing"
	"log"
)

func TestExecutePlugin(t *testing.T) {
	configFile := "../../cli/configs/default/default.conf"
	config,err := ReadConfig(configFile)
	if err != nil {
		t.Fail()
	}

	plugin := config.PluginDir + "/app/" + config.AppPlugin
	result := ExecutePlugin(config, "app", plugin, "--action", "quiesce")

	log.Println(result)

	if result.Code != 0 {
		t.Fail()
	}
}

func TestExecutePluginSimple(t *testing.T) {
	configFile := "../../cli/configs/default/default.conf"
	config,err := ReadConfig(configFile)
	if err != nil {
		t.Fail()
	}

	plugin := config.PluginDir + "/storage/" + config.StoragePlugin
	result := ExecutePluginSimple(config, "storage", plugin, "--action", "backup")

	log.Println(result)

	if result.Code != 0 {
		t.Fail()
	}
}