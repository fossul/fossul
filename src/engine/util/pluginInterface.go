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
package util

import (
	"errors"
	"fmt"
	"plugin"
	"regexp"
)

type AppPlugin interface {
	SetEnv(Config) Result
	Quiesce(Config) Result
	Unquiesce(Config) Result
	PreRestore(Config) Result
	PostRestore(Config) Result
	Discover(Config) DiscoverResult
	Info() Plugin
}

type StoragePlugin interface {
	SetEnv(Config) Result
	Backup(Config) Result
	Restore(Config) Result
	BackupDelete(Config) Result
	BackupDeleteWorkflow(Config) Result
	Mount(Config) Result
	Unmount(Config) Result
	BackupList(Config) Backups
	Info() Plugin
}

type ArchivePlugin interface {
	SetEnv(Config) Result
	Archive(Config) Result
	ArchiveDelete(Config) Result
	ArchiveList(Config) Archives
	Info() Plugin
}

func GetPluginPath(pluginName, pluginType string) string {
	var path string

	re := regexp.MustCompile(`\S+.so`)
	match := re.FindStringSubmatch(pluginName)

	if match != nil {
		path = "./plugins/" + pluginType + "/" + pluginName
	} else {
		fmt.Println("Native plugin [" + pluginName + "] does not exist, executing as basic plugin")
		path = ""
	}

	return path
}

func GetAppInterface(path string) (AppPlugin, error) {
	plugin, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	symPlugin, err := plugin.Lookup("AppPlugin")
	if err != nil {
		return nil, err
	}

	var appPlugin AppPlugin
	appPlugin, ok := symPlugin.(AppPlugin)
	if !ok {
		return nil, errors.New("Unexpected symbol type from module [ " + path + "], ensure plugin properly implements interface AppPlugin")
	}

	return appPlugin, nil
}

func GetStorageInterface(path string) (StoragePlugin, error) {
	plugin, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	symPlugin, err := plugin.Lookup("StoragePlugin")
	if err != nil {
		return nil, err
	}

	var storagePlugin StoragePlugin
	storagePlugin, ok := symPlugin.(StoragePlugin)
	if !ok {
		return nil, errors.New("Unexpected symbol type from module [ " + path + "], ensure plugin properly implements interface StoragePlugin")
	}

	return storagePlugin, nil
}

func GetArchiveInterface(path string) (ArchivePlugin, error) {
	plugin, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	symPlugin, err := plugin.Lookup("ArchivePlugin")
	if err != nil {
		return nil, err
	}

	var archivePlugin ArchivePlugin
	archivePlugin, ok := symPlugin.(ArchivePlugin)
	if !ok {
		return nil, errors.New("Unexpected symbol type from module [ " + path + "], ensure plugin properly implements interface ArchivePlugin")
	}

	return archivePlugin, nil
}
