package util

import (
	"fmt"
	"plugin"
)

type AppPlugin interface {
	SetEnv(Config) (Result)
	Quiesce() (Result)
	Unquiesce() (Result)
	Info() (Plugin)
}

type StoragePlugin interface {
	SetEnv(Config) (Result)
	Backup() (Result)
	BackupDelete() (Result)
	BackupList() ([]Backup)
	Info() (Plugin)
}

func GetPluginPath(pluginName string) string {
	var path string
	switch pluginName {
	case "mariadb.so":
		path = "./plugins/app/mariadb.so"
	case "container-basic.so":
		path = "./plugins/storage/container-basic.so"		
	default:
		fmt.Println("Built-in plugin [" + pluginName + "] does not exist")
		path = ""
	}

	return path
}

func GetAppInterface(path string) AppPlugin {
	plugin, err := plugin.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	symPlugin, err := plugin.Lookup("AppPlugin")
	if err != nil {
		fmt.Println(err)
	}

	var appPlugin AppPlugin
	appPlugin, ok := symPlugin.(AppPlugin)
	if !ok {
		fmt.Println(appPlugin,ok,"unexpected type from module symbol")
	}

	return appPlugin
}

func GetStorageInterface(path string) StoragePlugin {
	plugin, err := plugin.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	symPlugin, err := plugin.Lookup("StoragePlugin")
	if err != nil {
		fmt.Println(err)
	}

	var storagePlugin StoragePlugin
	storagePlugin, ok := symPlugin.(StoragePlugin)
	if !ok {
		fmt.Println(storagePlugin,ok,"unexpected type from module symbol")
	}

	return storagePlugin
}