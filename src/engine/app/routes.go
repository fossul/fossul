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
		"PreQuiesceCmd",
		"POST",
		"/preQuiesceCmd",
		PreQuiesceCmd,
	},
	Route{
		"QuiesceCmd",
		"POST",
		"/quiesceCmd",
		QuiesceCmd,
	},
	Route{
		"Quiesce",
		"POST",
		"/quiesce",
		Quiesce,
	},
	Route{
		"PostQuiesceCmd",
		"POST",
		"/postQuiesceCmd",
		PostQuiesceCmd,
	},
	Route{
		"PreUnquiesceCmd",
		"POST",
		"/preUnquiesceCmd",
		PreUnquiesceCmd,
	},
	Route{
		"UnquiesceCmd",
		"POST",
		"/unquiesceCmd",
		UnquiesceCmd,
	},
	Route{
		"Unquiesce",
		"POST",
		"/unquiesce",
		Unquiesce,
	},
	Route{
		"PostUnquiesceCmd",
		"POST",
		"/postUnquiesceCmd",
		PostUnquiesceCmd,
	},
	Route{
		"Discover",
		"POST",
		"/discover",
		Discover,
	},
	Route{
		"PreRestore",
		"POST",
		"/preRestore",
		PreRestore,
	},
	Route{
		"PostRestore",
		"POST",
		"/postRestore",
		PostRestore,
	},
	Route{
		"PreAppRestoreCmd",
		"POST",
		"/preAppRestoreCmd",
		PreAppRestoreCmd,
	},
	Route{
		"PostAppRestoreCmd",
		"POST",
		"/postAppRestoreCmd",
		PostAppRestoreCmd,
	},
}
