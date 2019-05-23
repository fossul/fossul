package util

import (
	"fmt"
	"plugin"
	"errors"
)

type AppPlugin interface {
	SetEnv(Config) (Result)
	Quiesce(Config) (Result)
	Unquiesce(Config) (Result)
	PreRestore(Config) (Result)
	PostRestore(Config) (Result)
	Discover(Config) (DiscoverResult)
	Info() (Plugin)
}

type StoragePlugin interface {
	SetEnv(Config) (Result)
	Backup(Config) (Result)
	Restore(Config) (Result)
	BackupDelete(Config) (Result)
	BackupList(Config) (Backups)
	Info() (Plugin)
}

type ArchivePlugin interface {
	SetEnv(Config) (Result)
	Archive(Config) (Result)
	ArchiveDelete(Config) (Result)
	ArchiveList(Config) ([]Archive)
	Info() (Plugin)
}

func GetPluginPath(pluginName string) string {
	var path string
	switch pluginName {
	case "mariadb.so":
		path = "./plugins/app/mariadb.so"
	case "mariadb-dump.so":
		path = "./plugins/app/mariadb-dump.so"	
	case "postgres.so":
		path = "./plugins/app/postgres.so"
	case "postgres-dump.so":
		path = "./plugins/app/postgres-dump.so"					
	case "mongo.so":
		path = "./plugins/app/mongo.so"
	case "mongo-dump.so":
		path = "./plugins/app/mongo-dump.so"				
	case "container-basic.so":
		path = "./plugins/storage/container-basic.so"		
	case "sample-app.so":
		path = "./plugins/app/sample-app.so"	
	case "sample-storage.so":
		path = "./plugins/storage/sample-storage.so"	
	case "sample-archive.so":
		path = "./plugins/archive/sample-archive.so"						
	default:
		fmt.Println("Native plugin [" + pluginName + "] does not exist, executing as basic plugin")
		path = ""
	}

	return path
}

func GetAppInterface(path string) (AppPlugin,error) {
	plugin, err := plugin.Open(path)
	if err != nil {
		return nil,err
	}

	symPlugin, err := plugin.Lookup("AppPlugin")
	if err != nil {
		return nil,err
	}

	var appPlugin AppPlugin
	appPlugin, ok := symPlugin.(AppPlugin)
	if !ok {
		return nil,errors.New("Unexpected symbol type from module [ " + path + "], ensure plugin properly implements interface AppPlugin")
	}

	return appPlugin,nil
}

func GetStorageInterface(path string) (StoragePlugin,error) {
	plugin, err := plugin.Open(path)
	if err != nil {
		return nil,err
	}

	symPlugin, err := plugin.Lookup("StoragePlugin")
	if err != nil {
		return nil,err
	}

	var storagePlugin StoragePlugin
	storagePlugin, ok := symPlugin.(StoragePlugin)
	if !ok {
		return nil,errors.New("Unexpected symbol type from module [ " + path + "], ensure plugin properly implements interface StoragePlugin")
	}

	return storagePlugin,nil
}

func GetArchiveInterface(path string) (ArchivePlugin,error) {
	plugin, err := plugin.Open(path)
	if err != nil {
		return nil,err
	}

	symPlugin, err := plugin.Lookup("ArchivePlugin")
	if err != nil {
		return nil,err
	}

	var archivePlugin ArchivePlugin
	archivePlugin, ok := symPlugin.(ArchivePlugin)
	if !ok {
		return nil,errors.New("Unexpected symbol type from module [ " + path + "], ensure plugin properly implements interface ArchivePlugin")
	}

	return archivePlugin,nil
}