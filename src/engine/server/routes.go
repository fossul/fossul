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
	"net/http"

	"github.com/fossul/fossul/src/engine/util"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = util.LogApi(handler, route.Name)

		if route.Name != "GetStatus" {
			handler = util.BasicAuth(handler, myUser, myPass)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"GetStatus",
		"GET",
		"/status",
		GetStatus,
	},
	Route{
		"StartBackupWorkflowLocalConfig",
		"GET",
		"/startBackupWorkflowLocalConfig",
		StartBackupWorkflowLocalConfig,
	},
	Route{
		"StartBackupWorkflow",
		"GET",
		"/startBackupWorkflow/{profileName}/{configName}/{policy}",
		StartBackupWorkflow,
	},
	Route{
		"StartOperatorBackupWorkflow",
		"GET",
		"/startOperatorBackupWorkflow/{profileName}/{configName}/{policy}",
		StartOperatorBackupWorkflow,
	},
	Route{
		"StartRestoreWorkflowLocalConfig",
		"GET",
		"/startRestoreWorkflowLocalConfig/{workflowId}",
		StartRestoreWorkflowLocalConfig,
	},
	Route{
		"StartRestoreWorkflow",
		"GET",
		"/startRestoreWorkflow/{profileName}/{configName}/{policy}/{workflowId}",
		StartRestoreWorkflow,
	},
	Route{
		"SendTrapSuccessCmd",
		"POST",
		"/sendTrapSuccessCmd",
		SendTrapSuccessCmd,
	},
	Route{
		"SendTrapErrorCmd",
		"POST",
		"/sendTrapErrorCmd",
		SendTrapErrorCmd,
	},
	Route{
		"GetConfig",
		"GET",
		"/getConfig/{profileName}/{configName}",
		GetConfig,
	},
	Route{
		"GetPluginConfig",
		"GET",
		"/getPluginConfig/{profileName}/{configName}/{pluginName}",
		GetPluginConfig,
	},
	Route{
		"GetDefaultConfig",
		"GET",
		"/getDefaultConfig",
		GetDefaultConfig,
	},
	Route{
		"GetDefaultPluginConfig",
		"GET",
		"/getDefaultPluginConfig/{pluginName}",
		GetDefaultPluginConfig,
	},
	Route{
		"GetWorkflowStepResults",
		"GET",
		"/getWorkflowStepResults/{profileName}/{configName}/{workflowId}/{stepId}",
		GetWorkflowStepResults,
	},
	Route{
		"GetWorkflowStatus",
		"GET",
		"/getWorkflowStatus/{profileName}/{configName}/{id}",
		GetWorkflowStatus,
	},
	Route{
		"DeleteWorkflowResults",
		"GET",
		"/deleteWorkflowResults/{profileName}/{configName}/{id}",
		DeleteWorkflowResults,
	},
	Route{
		"GetJobs",
		"GET",
		"/getJobs/{profileName}/{configName}",
		GetJobs,
	},
	Route{
		"AddConfig",
		"POST",
		"/addConfig/{profileName}/{configName}",
		AddConfig,
	},
	Route{
		"AddPluginConfig",
		"POST",
		"/addPluginConfig/{profileName}/{configName}/{pluginName}",
		AddPluginConfig,
	},
	Route{
		"DeleteConfig",
		"GET",
		"/deleteConfig/{profileName}/{configName}",
		DeleteConfig,
	},
	Route{
		"DeleteConfigDir",
		"GET",
		"/deleteConfigDir/{profileName}/{configName}",
		DeleteConfigDir,
	},
	Route{
		"DeletePluginConfig",
		"GET",
		"/deletePluginConfig/{profileName}/{configName}/{pluginName}",
		DeletePluginConfig,
	},
	Route{
		"AddProfile",
		"GET",
		"/addProfile/{profileName}",
		AddProfile,
	},
	Route{
		"DeleteProfile",
		"GET",
		"/deleteProfile/{profileName}",
		DeleteProfile,
	},
	Route{
		"ListProfiles",
		"GET",
		"/listProfiles",
		ListProfiles,
	},
	Route{
		"ListConfigs",
		"GET",
		"/listConfigs/{profileName}",
		ListConfigs,
	},
	Route{
		"ListPluginConfigs",
		"GET",
		"/listPluginConfigs/{profileName}/{configName}",
		ListPluginConfigs,
	},
	Route{
		"AddSchedule",
		"GET",
		"/addSchedule/{profileName}/{configName}/{policy}",
		AddSchedule,
	},
	Route{
		"DeleteSchedule",
		"GET",
		"/deleteSchedule/{profileName}/{configName}/{policy}",
		DeleteSchedule,
	},
	Route{
		"ListSchedules",
		"GET",
		"/listSchedules",
		ListSchedules,
	},
	Route{
		"GetSchedule",
		"GET",
		"/getSchedule/{profileName}/{configName}/{policy}",
		GetSchedule,
	},
	Route{
		"DeleteBackup",
		"GET",
		"/deleteBackup/{profileName}/{configName}/{policy}/{workflowId}",
		DeleteBackup,
	},
	Route{
		"CreateBackupCustomResource",
		"GET",
		"/createBackupCustomResource/{profileName}/{configName}/{policy}/{workflowId}/{timestamp}",
		CreateBackupCustomResource,
	},
	Route{
		"DeleteBackupCustomResource",
		"GET",
		"/deleteBackupCustomResource/{profileName}/{configName}/{policy}/{crName}",
		DeleteBackupCustomResource,
	},
	Route{
		"UpdateBackupCustomResource",
		"GET",
		"/updateBackupCustomResource/{profileName}/{configName}/{policy}/{crName}/{op}/{specKey}/{specValue}",
		UpdateBackupCustomResource,
	},
	Route{
		"BackupCustomResourceRetention",
		"GET",
		"/backupCustomResourceRetention/{profileName}/{configName}/{policy}",
		BackupCustomResourceRetention,
	},
	Route{
		"GetBackup",
		"GET",
		"/getBackup/{profileName}/{configName}/{workflowId}",
		GetBackup,
	},
}
