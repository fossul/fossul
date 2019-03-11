package main

import (
	"github.com/gorilla/mux"
    "net/http"
    "engine/util"
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
        "POST",
        "/pluginList",
        PluginList,
    }, 
    Route{
        "PluginInfo",
        "POST",
        "/pluginInfo/{plugin}",
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
        Backup,
    },
    Route{
        "BackupDelete",
        "POST",
        "/backupDelete",
        BackupDelete,
    },
}