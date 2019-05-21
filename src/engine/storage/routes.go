package main

import (
	"github.com/gorilla/mux"
    "net/http"
    "fossil/src/engine/util"
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