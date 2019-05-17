package util

import (
	"testing"
	"log"
	"os"
)

var pluginDir string = os.Getenv("FOSSIL_BUILD_PLUGIN_DIR")

func TestExecutePlugin(t *testing.T) {
	configFile := "../../cli/configs/default/default.conf"
	config,err := ReadConfig(configFile)
	if err != nil {
		t.Fail()
	}

	plugin := pluginDir + "/app/" + config.AppPlugin
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

	plugin := pluginDir + "/storage/" + config.StoragePlugin
	result := ExecutePluginSimple(config, "storage", plugin, "--action", "backup")

	log.Println(result)

	if result.Code != 0 {
		t.Fail()
	}
}