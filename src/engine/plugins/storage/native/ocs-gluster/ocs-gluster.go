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
	"fossul/src/engine/client/k8s"
	"fossul/src/engine/util"
	"github.com/heketi/heketi/client/api/go-client"
	//"strings"
)

type storagePlugin string

var StoragePlugin storagePlugin

func (s storagePlugin) SetEnv(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	result = util.SetResult(resultCode, messages)

	return result
}

func (s storagePlugin) Backup(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "Performing container backup")
	messages = append(messages, msg)

	heketi := client.NewClient("http://heketi-storage-app-storage.pu-ose2.coe.muc.redhat.com", "admin", "veWzd+r8MSgWsXClhvyYFL5bcmcDxX6sW1FigOGoetU=")

	// List clusters

	list, err := heketi.VolumeList()
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	for _, volumeId := range list.Volumes {
		volumeInfo, err := heketi.VolumeInfo(volumeId)
		if err != nil {
			msg := util.SetMessage("ERROR", err.Error())
			messages = append(messages, msg)

			result = util.SetResult(1, messages)
			return result
		}

		fmt.Println("HERE", volumeInfo.Name, volumeInfo.Mount, volumeInfo.Snapshot, volumeInfo.Size)
	}

	//output := strings.Join(list.Volumes, "\n")
	//fmt.Println("Clusters:", output)

	pvName, err := k8s.GetPersistentVolumeName(config.StoragePluginParameters["DatabaseNamespace"], config.StoragePluginParameters["PvcName"], config.StoragePluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	pvMountPath, err := k8s.GetGlusterPersistentVolumePath(pvName, config.StoragePluginParameters["AccessWithinCluster"])
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	fmt.Println("HERE1223", pvName, pvMountPath)

	podName, err := k8s.GetPodByName(config.StoragePluginParameters["Namespace"], config.StoragePluginParameters["PodName"], config.StoragePluginParameters["AccessWithinCluster"])
	fmt.Println("HERE1213", podName)
	if err != nil {
		msg := util.SetMessage("ERROR", err.Error())
		messages = append(messages, msg)

		result = util.SetResult(1, messages)
		return result
	}

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) Restore(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Restore ***")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) BackupDelete(config util.Config) util.Result {
	var result util.Result
	var messages []util.Message
	var resultCode int = 0

	msg := util.SetMessage("INFO", "*** Backup Delete ***")
	messages = append(messages, msg)

	result = util.SetResult(resultCode, messages)
	return result
}

func (s storagePlugin) BackupList(config util.Config) util.Backups {
	var backups util.Backups

	return backups
}

func (s storagePlugin) Info() util.Plugin {
	var plugin util.Plugin = setPlugin()
	return plugin
}

func setPlugin() (plugin util.Plugin) {
	plugin.Name = "ocs-gluster"
	plugin.Description = "OpenShift Container Storage (Gluster)"
	plugin.Version = "1.0.0"
	plugin.Type = "storage"

	var capabilities []util.Capability
	var backupCap util.Capability
	backupCap.Name = "backup"

	var backupListCap util.Capability
	backupListCap.Name = "backupList"

	var backupDeleteCap util.Capability
	backupDeleteCap.Name = "backupDelete"

	var infoCap util.Capability
	infoCap.Name = "info"

	capabilities = append(capabilities, backupCap, backupListCap, backupDeleteCap, infoCap)

	plugin.Capabilities = capabilities

	return plugin
}

func main() {}
