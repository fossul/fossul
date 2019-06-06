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
	"fossul/src/engine/util"
	"github.com/gorilla/mux"
	"net/http"
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
		"PluginList",
		"GET",
		"/pluginList/{pluginType}",
		PluginList,
	},
	Route{
		"PluginInfo",
		"POST",
		"/pluginInfo/{pluginName}/{pluginType}",
		PluginInfo,
	},
	Route{
		"Backup",
		"POST",
		"/backup",
		Backup,
	},
	Route{
		"BackupList",
		"POST",
		"/backupList",
		BackupList,
	},
	Route{
		"BackupDelete",
		"POST",
		"/backupDelete",
		BackupDelete,
	},
	Route{
		"BackupCreateCmd",
		"POST",
		"/backupCreateCmd",
		BackupCreateCmd,
	},
	Route{
		"BackupDeleteCmd",
		"POST",
		"/backupDeleteCmd",
		BackupDeleteCmd,
	},
	Route{
		"Archive",
		"POST",
		"/archive",
		Archive,
	},
	Route{
		"ArchiveList",
		"POST",
		"/archiveList",
		ArchiveList,
	},
	Route{
		"ArchiveDelete",
		"POST",
		"/archiveDelete",
		ArchiveDelete,
	},
	Route{
		"ArchiveCreateCmd",
		"POST",
		"/archiveCreateCmd",
		ArchiveCreateCmd,
	},
	Route{
		"ArchiveDeleteCmd",
		"POST",
		"/archiveDeleteCmd",
		ArchiveDeleteCmd,
	},
	Route{
		"Restore",
		"POST",
		"/restore",
		Restore,
	},
	Route{
		"RestoreCmd",
		"POST",
		"/restoreCmd",
		RestoreCmd,
	},
}
